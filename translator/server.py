import grpc
from concurrent import futures
import translation_pb2
import translation_pb2_grpc
from transformers import pipeline

class TranslatorService(translation_pb2_grpc.TranslatorServicer):
    def __init__(self):
        print("Loading model... this might take a moment.")
        # Load the model ONCE when the server starts
        self.pipe = pipeline("translation", model="Helsinki-NLP/opus-mt-de-es")
        print("Model loaded!")

    def Translate(self, request, context):
        # The actual translation logic
        print(f"Received request: {request.text}")
        
        # Run the translation
        result = self.pipe(request.text)
        translated_text = result[0]['translation_text']
        
        # Return the response defined in the proto
        return translation_pb2.TranslateResponse(translated_text=translated_text)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    translation_pb2_grpc.add_TranslatorServicer_to_server(TranslatorService(), server)
    
    # Listen on port 50051
    server.add_insecure_port('[::]:50051')
    print("Server starting on port 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
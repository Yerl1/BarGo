from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import joblib
from tensorflow.keras.models import load_model
import numpy as np
import os

app = FastAPI(title="Toxicity Detection API", version="1.0.0")

# Global variables for model and vectorizer
model = None
vectorizer = None

class TextRequest(BaseModel):
    comment: str

class PredictionResponse(BaseModel):
    toxic: bool

def load_ml_models():
    """Load ML model and vectorizer"""
    global model, vectorizer
    
    try:
        # Load vectorizer
        vectorizer = joblib.load('tfidf_vectorizer.joblib')
        print("✓ Vectorizer loaded successfully")
        
        # Load model
        model = load_model('toxic_comment_model.h5')
        print("✓ Model loaded successfully")
        
    except Exception as e:
        print(f"✗ Error loading models: {e}")
        raise e

@app.on_event("startup")
async def startup_event():
    """Load models on startup"""
    print("Starting up ML Service...")
    load_ml_models()
    print("ML Service ready!")

@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy", 
        "service": "toxicity-detector",
        "model_loaded": model is not None,
        "vectorizer_loaded": vectorizer is not None
    }

@app.post("/predict", response_model=PredictionResponse)
async def predict_toxicity(request: TextRequest):
    """Predict if text is toxic"""
    try:
        if model is None or vectorizer is None:
            raise HTTPException(status_code=503, detail="Models not loaded")
        
        # Transform text using TF-IDF
        text_tfidf = vectorizer.transform([request.comment])
        
        # Predict
        prediction = model.predict(text_tfidf.toarray())
        is_toxic = bool(prediction[0][0] > 0.5)
        
        return PredictionResponse(toxic=is_toxic)
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Prediction error: {str(e)}")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=3005)
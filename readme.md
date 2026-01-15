# AI Bill Splitter

Split bills intelligently using AI and natural language processing. Upload a bill image, define who paid for what with simple rules, and get instant fair calculations—no more awkward money conversations at dinner, no more time wastage after shopping just calculate bill share.

## Features

✨ **AI-Powered Image Recognition** – Upload bill photos and automatically extract items and prices  
🤖 **Natural Language Rules** – Define sharing rules in plain English (e.g., "Alice paid for appetizers")  
⚖️ **Accurate Calculations** – Handles tax distribution and complex sharing scenarios  
🎯 **Fair & Transparent** – See exactly how each person's share breaks down by item  
🚀 **Fast & Reliable** – Built with Go for performance; powered by Google's Gemini API  

## How It Works

1. **Upload a Bill** – Take a photo or upload an image of your receipt
2. **Define Sharing Rules** – Describe how items should be split (or let everyone split everything equally)
3. **Get Results** – Instantly see who owes what, broken down by item

### Example

**Bill:**
| Item | Price | Qty | Tax | Total |
|------|-------|-----|-----|-------|
| Appetizer | $50 | 1 | $5 | $55 |
| Entree | $30 | 2 | $6 | $66 |
| Dessert | $20 | 1 | $2 | $22 |

**Rules:**
- Appetizer shared by Alice and Bob
- Entrees split equally among Alice, Bob, and Charlie
- Charlie had the dessert

**Result:**
| Person | Appetizer | Entrees | Dessert | Total |
|--------|-----------|---------|---------|-------|
| Alice | $27.50 | $22.00 | — | $49.50 |
| Bob | $27.50 | $22.00 | — | $49.50 |
| Charlie | — | $22.00 | $22.00 | $44.00 |

## Tech Stack

- **Backend:** Go (Golang)
- **Frontend:** HTML/CSS/JavaScript (Entirely Vibe Coded 🙂)
- **AI Model:** Google Gemini 2.5 Flash
- **Image Processing:** OCR via Gemini API
- **Deployment:** Docker, Google Cloud Run

## Architecture

- **main.go** – Server setup and routing
- **http-handlers.go** – API endpoints
- **query-model.go** – AI model interaction
- **process-bill.go** – Bill parsing
- **calculate-split.go** – Core splitting logic
- **validations.go** – Model output validations
- **data-models.go** – Data structures

## API Endpoints

- `GET /` – Web interface
- `POST /split` – Process bill image and rules (returns calculated splits)
- `GET /hi` – Health check

### Project Status

✅ Core splitting logic  
✅ AI image recognition  
✅ Web interface  
✅ API endpoints  
⏳ Deployment (in progress)

## Future Enhancements

- Split distribution sharing
- Support for multiple currencies
- User authentication and bill history
- Receipt OCR improvements
- Group expense tracking

Built with ❤️ for fairness in shared expenses.

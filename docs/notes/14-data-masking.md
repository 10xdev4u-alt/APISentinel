# 14: Masking Sensitive Data (DLP Level 2) 🎭

Blocking a response because it contains a single credit card number is the safest approach, but sometimes it ruins the User Experience. 

## The Balanced Approach: Masking
Instead of returning a `500 Internal Server Error`, we can **modify** the response body on the fly. We find the sensitive string and replace it with something safe.

### Example:
- **Original:** `Order processed for card 4111-1111-1111-1234.`
- **Masked:**   `Order processed for card ****-****-****-1234.`

## How it works in API Sentinel
1.  **Capture:** We already capture the response body in our `DLPMiddleware`.
2.  **Inspect:** We run our RegEx.
3.  **Replace:** If we find a match and the config says `action: mask`, we use `Regexp.ReplaceAllStringFunc` to obfuscate the data.
4.  **Forward:** We send the modified body to the user.

## Why use Masking?
- **Internal Tools:** You might want to see that *a* card was used without seeing the *actual* card.
- **Support Logs:** Masking allows developers to debug issues without being exposed to PCI/PII data.

Next, we will update our DLP logic!

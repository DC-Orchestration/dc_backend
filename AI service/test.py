import json
import difflib
from typing import Dict, List, Any
import re
from langchain_groq import ChatGroq
from langchain.schema import HumanMessage

class TextComparator:
    def __init__(self, model_name: str = "llama3-70b-8192"):
        self.model_name = model_name
        self.chat = ChatGroq(
            groq_api_key="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",
            model_name=model_name
        )

    async def get_line_differences(self, original: str, modified: str) -> List[str]:
        """Get line-by-line differences using difflib"""
        diff = difflib.unified_diff(
            original.splitlines(),
            modified.splitlines(),
            lineterm='',
            fromfile='original',
            tofile='modified'
        )
        return list(diff)
        
    async def analyze_differences(self, original_text: str, modified_text: str) -> Dict:
        """Use LLM to categorize and describe differences"""
        # Prepare text snippets safely
        orig_snippet = original_text[:2000].replace('\n', '\\n') + ("..." if len(original_text) > 2000 else "")
        mod_snippet = modified_text[:2000].replace('\n', '\\n') + ("..." if len(modified_text) > 2000 else "")
        diff_lines = await self.get_line_differences(original_text, modified_text)
        diff_snippet = '\n'.join(diff_lines[:50]).replace('\n', '\\n') + ("..." if len(diff_lines) > 50 else "")

        prompt = f"""
        Analyze these text differences and return a plain text response (ONLY the response, no markdown or extra text):

        ORIGINAL TEXT:
        {orig_snippet}

        MODIFIED TEXT:
        {mod_snippet}

        DIFFERENCES (line-by-line):
        {diff_snippet}

        Generate the summary output of the changes with the following rules in mind:
        - Focus on structural changes first
        - Mention line number ranges for major changes
        - Classify as: Added Content/Formatting/Refactor/Text Edit/Removed Content
        - Maximum 3 sentences for summary

        Return ONLY the summary as plain text with no additional text or markdown formatting.
        """

        response = await self.invoke_model(prompt)
        print("LLM Response:", response)
        # Remove excessive newlines
        return response

    async def compare_texts(self, original_text: str, modified_text: str) -> Dict:
        """Main comparison workflow for direct text input"""
        try:
            # Get raw differences
            #diff_lines = await self.get_line_differences(original_text, modified_text)

            # Get analyzed differences from LLM
            analyzed_diff = await self.analyze_differences(original_text, modified_text)

            # Prepare final output
            return analyzed_diff

        except Exception as e:
            print(f"Comparison failed: {e}")
            return {"error": str(e)}

    async def invoke_model(self, prompt: str) -> str:
        """Wrapper for LLM calls"""
        try:
            response = self.chat.invoke([HumanMessage(content=prompt)])
            return response.content.strip()
        except Exception as e:
            print(f"Error invoking model: {e}")
            raise

    def extract_valid_json(self, response: str) -> Dict:
        """Extract JSON from LLM response"""
        try:
            # Try direct parse first
            return json.loads(response)
        except json.JSONDecodeError:
            # Fallback: find first {...} in response
            match = re.search(r'\{.*\}', response, re.DOTALL)
            if match:
                try:
                    return json.loads(match.group(0))
                except json.JSONDecodeError as e:
                    print(f"Failed to parse JSON: {e}")
        return {"error": "Could not extract valid JSON"}

import asyncio
async def test_comparator():
    """Test function for TextComparator"""
    print("\n=== Testing TextComparator ===")
    comparator = TextComparator()
    
    # Test Case 1: Simple text change
    print("\nTest 1: Simple text modification")
    result1 = await comparator.compare_texts(
        f"Hello world! this is a test. This is a new line.",
        f"Hello world! this is a test. This is an edited line."
    )
    print("Result 1:", result1)
    
    # Test Case 3: No changes
    print("\nTest 3: Identical texts")
    result3 = await comparator.compare_texts(
        "Identical content",
        "Identical content"
    )
    print("Result 2:", result3)

if __name__ == "__main__":
    # Run tests when executed directly
    asyncio.run(test_comparator())
    
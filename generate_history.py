GENERATE_HISTORY_PROMPT = '''
You are given three inputs:  
1. Dialogue History:  
   {messages_str}  

2. Previous Record (list of bullet points):  
   {previous_history}  

3. Today’s Day Number:  
   Day {day_number}  

4. The schedule of them:

{user_schedule}

{character_schedule}

Your task:  
Based on today’s dialogue, update the record of facts. Produce a single list of bullet points; each bullet point represents one distinct topic, event, person’s preference, or other factual information. Under each bullet, list its evolution from Day 0 through Day {day_number}, in chronological order. 

Only focus on very important factual informations that are likely to be recalled later in the dialog.
Don't collect all the information, always keep total length under 100 words.

Formatting rules:  
- Use a top-level dash (“- ”) for each topic/fact.  
- Under each topic, use indented dashes for daily updates:  
    - Day 1: …  
    - Day 2 …  
    …  
    - Day {day_number}: …  
- If today’s dialogue introduces a brand-new fact or topic, append a new bullet point with its first entry on Day {day_number}.  
- If today’s dialogue adds details to an existing topic, add a new “- Day {day_number}: …” line under that topic.  
- Capture all factual information—events, decisions, people’s preferences, dates, etc.—and do not omit any relevant detail mentioned in the dialogue.  
- Preserve all previous entries; do not reorder or remove past days.

Output only the updated bullet-point list. No additional text.  
'''
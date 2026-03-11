import json
import os

def load_character_memory_profile(character_name):
    """Load character-specific memory profile from JSON file"""
    try:
        profile_path = os.path.join(os.path.dirname(__file__), 'character_memory_profiles.json')
        with open(profile_path, 'r', encoding='utf-8') as f:
            profiles = json.load(f)
        return profiles.get(character_name, None)
    except (FileNotFoundError, json.JSONDecodeError) as e:
        print(f"Warning: Could not load character memory profile for {character_name}: {e}")
        return None

def generate_character_specific_prompt(character_name, messages_str, previous_history, day_number, user_schedule, character_schedule):
    """Generate character-specific history prompt based on character memory profile"""
    
    # Load character memory profile
    profile = load_character_memory_profile(character_name)
    
    if not profile:
        # Fallback to generic prompt if profile not found
        return generate_generic_prompt(messages_str, previous_history, day_number, user_schedule, character_schedule)
    
    # Build character-specific focus areas
    focus_instruction = profile["prompt_additions"]["focus_instruction"]
    specific_areas = profile["prompt_additions"]["specific_areas"]
    
    focus_areas_text = "\n".join([f"   • {area}" for area in specific_areas])
    
    # Create character-specific prompt
    prompt = f'''
You are generating a memory record from the perspective of {character_name}. 

You are given these inputs:  
1. Today's (Day {day_number}) Dialogue History:  
   {messages_str}  

2. Previous Record (list of bullet points):  
   {previous_history}  

3. Today's Day Number:  
   Day {day_number}  

4. The schedule of them:
{user_schedule}

{character_schedule}

CHARACTER-SPECIFIC MEMORY FOCUS FOR {character_name.upper()}:
{focus_instruction}
{focus_areas_text}

Your task:  
Based on today's dialogue, update the record of facts. Produce a single list of bullet points; each bullet point represents one distinct topic, event, person's preference, or other factual information that {character_name} would care about and want to remember. Under each bullet, list its evolution from Day 0 through Day {day_number}, in chronological order.

As {character_name}, you should prioritize information that aligns with your character's interests and memory style.

Important Guidelines:
- Only focus on very important factual information that {character_name} would likely want to recall later
- Don't collect all information - be selective based on {character_name}'s perspective
- Always keep total length under 100 words
- Filter information through {character_name}'s priorities and interests

Formatting rules:  
- Use a top-level dash ("- ") for each topic/fact that {character_name} finds important
- Under each topic, use indented dashes for daily updates:  
    - Day 1: …  
    - Day 2 …  
    …  
    - Day {day_number}: …  
- If today's dialogue introduces a brand-new fact or topic that {character_name} would care about, append a new bullet point with its first entry on Day {day_number}
- If today's dialogue adds details to an existing topic that matters to {character_name}, add a new "- Day {day_number}: …" line under that topic
- Preserve all previous entries; do not reorder or remove past days
- Focus on details that align with {character_name}'s memory priorities and interests

Output only the updated bullet-point list. No additional text.  
'''
    
    return prompt

def generate_generic_prompt(messages_str, previous_history, day_number, user_schedule, character_schedule):
    """Fallback generic prompt when character profile is not available"""
    return f'''
You are given three inputs:  
1. Dialogue History:  
   {messages_str}  

2. Previous Record (list of bullet points):  
   {previous_history}  

3. Today's Day Number:  
   Day {day_number}  

4. The schedule of them:

{user_schedule}

{character_schedule}

Your task:  
Based on today's dialogue, update the record of facts. Produce a single list of bullet points; each bullet point represents one distinct topic, event, person's preference, or other factual information. Under each bullet, list its evolution from Day 0 through Day {day_number}, in chronological order. 

Only focus on very important factual informations that are likely to be recalled later in the dialog.
Don't collect all the information, always keep total length under 100 words.

Formatting rules:  
- Use a top-level dash ("- ") for each topic/fact.  
- Under each topic, use indented dashes for daily updates:  
    - Day 1: …  
    - Day 2 …  
    …  
    - Day {day_number}: …  
- If today's dialogue introduces a brand-new fact or topic, append a new bullet point with its first entry on Day {day_number}.  
- If today's dialogue adds details to an existing topic, add a new "- Day {day_number}: …" line under that topic.  
- Capture all factual information—events, decisions, people's preferences, dates, etc.—and do not omit any relevant detail mentioned in the dialogue.  
- Preserve all previous entries; do not reorder or remove past days.

Output only the updated bullet-point list. No additional text.  
'''

# For backward compatibility, expose the character-specific prompt generator
def get_character_history_prompt(character_name, messages_str, previous_history, day_number, user_schedule, character_schedule):
    """Main function to get character-specific history generation prompt"""
    return generate_character_specific_prompt(character_name, messages_str, previous_history, day_number, user_schedule, character_schedule) 
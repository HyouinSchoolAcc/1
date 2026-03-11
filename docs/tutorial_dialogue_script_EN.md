# 📜 Divergence 2% Tutorial Dialogue Script
## *"Cao Cao's Guide to Writing Characters"*
### Version 2.0 (English) — January 2026

---

## 📋 Document Overview

**Purpose:** Gal-game style interactive tutorial for onboarding new writers  
**Guide Character:** Cao Cao (曹操)  
**Estimated Duration:** ~10 minutes  
**Tone:** Fun, irreverent, fourth-wall-breaking humor with helpful information

### Section Structure
1. Just Joining — Greetings
2. What Do We Do?
3. Navigation Page
4. Writing Interface
5. Character Descriptions
6. Writer's Lounge
7. Summary
8. Ending
9. Easter Eggs (Optional)

*Note: Creating your own character has a separate tutorial on the /newcharacter page.*

### Format Legend
- `Cao Cao:` — Guide character speaks
- `User:` — Player/user speaks
- `→ CHOICE:` — Branching choice point
- `[ HIGHLIGHT: xxx ]` — UI element gets highlighted
- `[ LOAD: url ]` — Page navigation in iframe
- `[ PRACTICE ]` — User interaction mode
- `→ [Action]` — Navigation or system action

---

## SECTION 1: Just Joining — Greetings

*First impressions and setting expectations*

\`\`\`
Cao Cao: Welcome! A warm welcome!

→ CHOICE:
  - "So enthusiastic? That doesn't seem like you." 
      → [Select fun mode, continue to 1.1]
  
  - "Just tell me what to do (Skip to business)" 
      → [Select business mode, jump to 1.2]
\`\`\`

### 1.1 Fun Path

\`\`\`
Cao Cao: We're short on writers, hehe.
Cao Cao: Honestly, I resisted doing PR, but the data quality forced my hand.

User: Forced your hand... is that the right expression?

Cao Cao: ...Cut the crap.
Cao Cao: Let's get straight to the point. Our current data needs you to play two roles.
Cao Cao: You play both "the person chatting with AI" and "the AI itself".

User: ...That sounds like writing a novel.

Cao Cao: No no no, novels have narration.
Cao Cao: We are pure chat logs.
Cao Cao: Just like venting to a friend—real, grounded, emotional.

→ [Continue to Section 2]
\`\`\`

### 1.2 Business Path

\`\`\`
Cao Cao: Okay, the short version:
Cao Cao: We need you to write dialogue data between AI characters and users.
Cao Cao: Play both roles, simulate real chat.
Cao Cao: Pass the quality check, get paid.

→ [Continue to Section 2]
\`\`\`

---

## SECTION 2: What Do We Do?

*Explaining the core workflow*

\`\`\`
Cao Cao: Let's officially talk about the process.
Cao Cao: Step 1: Understand the system, view character settings.
Cao Cao: Step 2: Create a "User" (the person chatting with the character).
Cao Cao: Step 3: Simulate a whole day's chat.
Cao Cao: Step 4: Submit for review, get paid.

User: That's it?

Cao Cao: That's it.
Cao Cao: But "natural dialogue" is harder than you think.
Cao Cao: Have you ever tried pretending to be someone else for 30 minutes?

User: No.

Cao Cao: Well, you're about to experience it.
Cao Cao: The core of this job is "immersion".
Cao Cao: You have to make the reviewer feel "This really looks like two people chatting".
\`\`\`

---

## SECTION 3: Navigation Page

*Introduction to the main interface*

\`\`\`
Cao Cao: Next, I'll show you the actual interface.
Cao Cao: Let's go to the character selection page first.

[ LOAD: /navigation/e ]

Cao Cao: This is the character selection page.
Cao Cao: Currently, there are three characters to write.
Cao Cao: Kurisu, Kafka, and Lin Lu.

Cao Cao: By the way, there is a Home button on the page.
Cao Cao: Click it to return home if you get lost.

[ HIGHLIGHT: .character-card (multiple arrows) ]

Cao Cao: Pick a character!

[Click any character card to continue]
\`\`\`

---

## SECTION 4: Writing Interface

*The main writing experience*

\`\`\`
[ LOAD: /kurisu/e ]

Cao Cao: Okay, this is the writing interface.
Cao Cao: First, you need to create a temporary character to experience it.
Cao Cao: Click the "Add Temporary Character" button in the left sidebar.
Cao Cao: Note: "What does this person want most right now?" is important. It determines the character's motivation.
Cao Cao: Go ahead and write your character, I'll wait.

[ WAIT: character creation ]

Cao Cao: Great! Your temporary character has been created.
Cao Cao: See the sidebar on the left? That's the dataset list.
Cao Cao: Each item represents a day's data. Click to open that day's dialogue.
Cao Cao: Now try clicking the character you just created.

[ HIGHLIGHT: sidebar item ]

[Click to continue]

Cao Cao: The right panel has several areas:
Cao Cao: Info Area — Shows info about the current data.
Cao Cao: Comments Area — Where other writers leave comments.
Cao Cao: Review Area — Feedback from editors after review.

Cao Cao: There are a few shortcuts for writing dialogue:
Cao Cao: Enter = New line (Same person continues messaging).
Cao Cao: Shift + Enter = Switch speaker.
Cao Cao: And the time bar—used to mark the timestamp of each message.
Cao Cao: On the first day, try to ask information regarding values, life styles and habits.
Cao Cao: Now try writing something in the input box.

[ HIGHLIGHT: textarea ]
[ PRACTICE: Use Enter for new line, Shift+Enter to switch speaker ]

Cao Cao: Got the hang of it?

User: Yeah.

Cao Cao: Remember to click the Save button!

[ HIGHLIGHT: save button ]

Cao Cao: ⚠️ Important Reminder:
Cao Cao: You just created a temporary character.
Cao Cao: If you like what you just wrote...
Cao Cao: Register now! Otherwise, content will be lost when you close the page.

[ REGISTER PROMPT ]
  - "📝 Sign Up to Save Your Work" → /login
  - "Continue without saving →" → continue
\`\`\`

---

## SECTION 5: Character Descriptions

*Understanding character settings*

\`\`\`
[ LOAD: /descriptions/e ]

Cao Cao: Next, let's look at character settings.
Cao Cao: Here are the detailed settings for each character.
Cao Cao: Take a look for a moment.

[ BROWSE PAUSE: 8 seconds ]

Cao Cao: Done? Let's talk about the key points.
Cao Cao: You must check two things before writing:
Cao Cao: "Background Settings" — Character's worldview and relationships.
Cao Cao: "Daily Settings" — What happened that day.
Cao Cao: These affect the direction and tone of the dialogue.

User: Got it, write based on the day's context.

Cao Cao: Right. Check back if unsure, don't make it up.
\`\`\`

---

## SECTION 6: Writer's Lounge

*Community and communication*

\`\`\`
[ LOAD: /lounge/e ]

Cao Cao: Finally, let's check the Writer's Lounge.
Cao Cao: This is the Writer's Lounge, where you can chat with other writers and the dev team.
Cao Cao: You can ask questions here.
Cao Cao: By the way, there are three buttons in the bottom left.
Cao Cao: Used for sharing content, exporting data, etc.

User: Understood.
\`\`\`

---

## SECTION 7: Summary

*Recap and fun facts*

\`\`\`
Cao Cao: Alright, that's the process.
Cao Cao: Summary:
Cao Cao: 1️⃣ Pick a character
Cao Cao: 2️⃣ Read settings (Very important!)
Cao Cao: 3️⃣ Create "User" identity
Cao Cao: 4️⃣ Write a full day's dialogue
Cao Cao: 5️⃣ Submit for review → Get paid!

Cao Cao: By the way, here are a few fun facts.
Cao Cao: Things you might not see in the docs.

Cao Cao: 【Kurisu】Secretly fantasizes about superpowers.
Cao Cao: Secretly thinks "If only I had telekinesis".

Cao Cao: 【Kafka】Has stalker tendencies, weighs "is it worth it".

User: That's intense.

Cao Cao: 【Lin Lu】Has the intellectual's "Curse of Knowledge".
Cao Cao: Doesn't understand why fixing a computer can't just be "restart it".

Cao Cao: Finally, a small quiz.

User: A test?!

Cao Cao: Don't worry, it's open book.

→ CHOICE:
  - "I'm ready, start the quiz" 
      → [Navigate to /quiz/e]
  
  - "Let me look at the settings again" 
      → [Navigate to /descriptions/e]
  
  - "Skip, start writing directly" 
      → [Jump to Ending]
\`\`\`

---

## SECTION 8: Ending

*Tutorial wrap-up and send-off*

\`\`\`
Cao Cao: Welcome to Singularity Gate, Researcher.
Cao Cao: Go pick a character and start.
Cao Cao: I'm on standby here.
Cao Cao: Even though I'm wearing this weird outfit...

User: (Laughs)

Cao Cao: Don't laugh!!

→ CHOICE:
  - "🖊️ Start Writing" 
      → [Navigate to /navigation/e]
  
  - "📖 Check Settings" 
      → [Navigate to /descriptions/e]

[ END OF TUTORIAL ]
\`\`\`

---

## SECTION 9: Easter Eggs (Optional Features)

### 9.1 Avatar Interaction
*If user clicks Cao Cao's avatar multiple times:*

\`\`\`
[Click 1]
Cao Cao: Why are you poking me?

[Click 2]
Cao Cao: ...

[Click 3]
Cao Cao: Poke me again and I'll kick you out.

[Click 4]
Cao Cao: ...

[Click 5]
Cao Cao: Okay okay, I know you're bored. Let's continue the tutorial.
\`\`\`

### 9.2 Idle Detection
*If user doesn't interact for 30+ seconds:*

\`\`\`
Cao Cao: ...Hey.
Cao Cao: Anyone there?
Cao Cao: Did you go to the bathroom or something?
Cao Cao: I'll wait...
\`\`\`

### 9.3 Quiz Failure Response
*If user fails the character quiz:*

\`\`\`
Cao Cao: ...
Cao Cao: Did you actually read the settings?
Cao Cao: Go back and look again. These questions were free points.
Cao Cao: When I lost at the Battle of Red Cliffs, even that wasn't this embarrassing.

User: You really want to bring up Red Cliffs?

Cao Cao: ...I take that back. You can try again.
\`\`\`

---

## 📝 Implementation Notes

### Key Features to Highlight in Tutorial
| Feature | Where to Mention | UI Element |
|---------|------------------|------------|
| Home button | Section 3 | \`.home-btn\` |
| 3 Characters | Section 3 | \`.character-card\` |
| Sidebar datasets | Section 4 | \`.day-item\` |
| Keyboard shortcuts | Section 4 | N/A (verbal) |
| Timestamp bar | Section 4 | Time input fields |
| Save button | Section 4 | \`.save-btn\` |
| Writers' Lounge | Section 6 | Lounge page |

### Story Elements Incorporated
| Story | Section | Notes |
|-------|---------|-------|
| Kurisu's superpower fantasies | Section 7 | Fun fact, embarrassed if caught |
| Kafka's stalking tendencies | Section 7 | Internal conflict angle |
| Lin Lu's intellect complex | Section 7 | "Curse of knowledge" framing |

### Branching Structure Summary
\`\`\`
START (Intro Screen)
  │
  ├─ "Let's Go!" → Section 1
  │
  └─ "Skip" → /navigation/e
  
Section 1
  │
  ├─ Fun Path → Section 1.1 → Section 2
  │
  └─ Business Path → Section 1.2 → Section 2
  
Section 2-6 (Linear with interactions)
  │
  ↓
Section 7 (Summary + Quiz Choice)
  │
  ├─ Start Quiz → /quiz/e
  ├─ Review Settings → /descriptions/e
  └─ Skip → Section 8
  
Section 8 (Ending)
  │
  ├─ Start Writing → /navigation/e
  └─ Check Settings → /descriptions/e
\`\`\`

---

## 🎮 Tone Guidelines

### Do's ✅
- Fourth-wall breaking humor
- Self-deprecating jokes (Cao Cao's outfit, Red Cliffs reference)
- Casual internet slang and memes
- Quick back-and-forth banter
- Genuine helpfulness underneath the jokes

### Don'ts ❌
- Don't be condescending
- Don't make the tutorial feel like a lecture
- Don't skip the fun parts for efficiency
- Don't break character (Cao Cao should always feel like Cao Cao)
- Don't make jokes at the expense of users

### Voice Reference
Cao Cao in this tutorial is:
- Slightly sarcastic but ultimately supportive
- Annoyed by his own costume
- Knowledgeable but not preachy
- Quick to banter with the user
- Secretly proud of the project

---

## 📚 Context for Non-Chinese Readers

### Who is Cao Cao?
Cao Cao (155–220 AD) was a legendary warlord during China's Three Kingdoms period. He's one of the most iconic figures in Chinese history and literature, known for being cunning, ambitious, and complex. In this tutorial, he's used as a humorous guide character dressed in traditional costume—a deliberate absurdist choice that the character himself finds ridiculous.

### Character References
- **Makise Kurisu**: Protagonist from the visual novel/anime *Steins;Gate*. An 18-year-old neuroscience genius known for her tsundere personality.
- **Kafka**: Character from the gacha game *Honkai: Star Rail*. A mysterious and elegant antagonist-turned-ally with a dangerous charm.
- **Lin Lu**: An original character created for this project. A gentle literature professor with an intellectual's blind spots.

### The Red Cliffs Reference
The Battle of Red Cliffs (208 AD) was Cao Cao's most famous defeat. It's frequently referenced in Chinese media as his greatest embarrassment. The joke in the quiz failure section plays on this historical fact.

---

*Document synchronized with guide.html*  
*Last updated: January 2026*

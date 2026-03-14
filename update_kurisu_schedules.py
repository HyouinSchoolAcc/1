import json
import shutil

# Read the production file
import os, tempfile

with open(r'C:\Users\user\Desktop\data_labeler\data\character_profiles.json', 'r', encoding='utf-8') as f:
    data = json.load(f)

# ============================================================
# NEW SCHEDULE ENTRIES FOR KURISU (Days 3-20)
# ============================================================

new_days = [
    {
        "day": 3,
        "title": "\u7b2c3\u5929\uff1a\u7cfb\u7edf\u8bbe\u8ba1",
        "title_en": "Day 3: System Design",
        "morning": "\u65e9\u4e0a\u5403\u4e86\u71d5\u9ea6\uff0c\u7136\u540e\u5f00\u59cb\u8bbe\u8ba1\u6574\u4e2a\u5b9e\u9a8c\u7cfb\u7edf\u3002\n\n\u7cfb\u7edf\u8bbe\u8ba1\u5206\u4e24\u4e2a\u9636\u6bb5\uff1a\n\n\u7b2c\u4e00\u9636\u6bb5\uff1a\u9a8c\u8bc1\u4eba\u5de5\u5927\u8111\u529f\u80fd\n\u642d\u5efa\u57fa\u7840\u7684\"\u7535\u6781\u6d4b\u8bd5\u5957\u4ef6\"\uff1a\n- \u6cd5\u62c9\u7b2c\u7b3c\n- \u57f9\u517b\u8111\u7ec6\u80de\n- \u642d\u5efa\u4fa7\u58c1\u7535\u6781\n- \u4ece\u7535\u6781\u5f15\u7ebf\u51fa\u6cd5\u62c9\u7b2c\u7b3c\u6765\u8bfb\u53d6\u4fe1\u53f7\n- \u8ba9\u5927\u8111\u5b66\u4f1a\u73a9\u300a\u6bc1\u706d\u6218\u58eb\u300b(Doom) \u6765\u9a8c\u8bc1\u529f\u80fd\u6b63\u5e38\n\n\u7b2c\u4e8c\u9636\u6bb5\uff1a\u5f15\u5165\u91cf\u5b50\u4fe1\u53f7\u9a8c\u8bc1\u5668\u6d4b\u8bd5\u7cfb\u7edf\n- \u8ba9\u9694\u58c1\u5b9e\u9a8c\u5ba4\u53d1\u5c04\u4fe1\u53f7\n- \u5728\u7bb1\u5b50\u672b\u7aef\u88c5\u5207\u4f26\u79d1\u592b\u63a2\u6d4b\u5668\u83b7\u53d6\u4e00\u7ef4\u4fe1\u53f7\n- \u9a8c\u8bc1\u8111\u7ec4\u7ec7\u5728\u8fd9\u79cd\u6761\u4ef6\u4e0b\u8fd8\u80fd\u4e0d\u80fd\u73a9Doom",
        "morning_en": "Ate oatmeal for breakfast, then started designing the entire experimental system.\n\nSystem design in two stages:\n\nStage 1: Verify artificial brain functions\nBuild the basic \"electrode testing kit\":\n- Faraday cage\n- Growing brain cells\n- Building the side wall electrodes\n- Connecting everything from the electrodes out through the Faraday cage for signal readout\n- Have the brain trained to play Doom to verify working function\n\nStage 2: Implement a quantum signal verifier to test the system\n- Have the neighboring lab shoot signals through\n- Cherenkov detectors at the end of the box for a 1-D signal\n- Test whether the brain tissues are still capable of playing Doom",
        "afternoon": "\u7ed5\u6821\u56ed\u8dd1\u4e86\u4e00\u5708\u3002\u4e0a\u4e86\u7b2c\u4e00\u8282\u91cf\u5b50\u573a\u8bba\u8bfe\u3002\u505a\u4e86\u751f\u7269\u5b9e\u9a8c\u5ba4\u7684\u7ebf\u4e0a\u57f9\u8bad\u3002",
        "afternoon_en": "Ran around campus. Took the first quantum field theory class. Completed bio lab online training.",
        "evening": "\u770b\u4e86\u4e00\u4e9b\u89c6\u9891\uff0c\u5b66\u4e60\u522b\u4eba\u662f\u600e\u4e48\u89e3\u51b3\u7c7b\u4f3c\u95ee\u9898\u7684\u3002",
        "evening_en": "Watched some videos on how others have tackled the problem.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 4,
        "title": "\u7b2c4\u5929\uff1a\u7ec6\u5316\u8bbe\u8ba1",
        "title_en": "Day 4: Fleshing Out Designs",
        "morning": "\u7ec6\u5316\u8bbe\u8ba1\u548c\u89c4\u683c\u53c2\u6570\u3002\u4e0a\u624bSolidWorks\uff0c\u7ed3\u679c\u53d1\u73b0\u7cfb\u7edf\u9700\u8981\u7528\u751f\u7269\u517c\u5bb9\u6750\u6599\uff0c\u5f00\u59cb\u5230\u5904\u627e\u65b0\u6750\u6599\u3002",
        "morning_en": "Fleshed out more of the designs and specifications. Started getting right into SolidWorks, then realized the system needs to be on bio-friendly materials\u2014started searching for new materials.",
        "afternoon": "\u53bb\u751f\u7269\u5b9e\u9a8c\u5ba4\u505a\u5b89\u5168\u57f9\u8bad\u3002\u8fd8\u770b\u5230\u4e86\u7528\u8367\u5149\u6807\u8bb0\u7684\u53d1\u5149\u7ec6\u80de\uff0c\u631a\u9177\u7684\u3002",
        "afternoon_en": "Went to the bio lab for safety training. Also got to see glowing cells with markers\u2014pretty cool.",
        "evening": "\u5c1d\u8bd5\u8ba2\u8d2d\u4e00\u4e9b\u8584\u819c\u6765\u505a\u7535\u6781\u3002",
        "evening_en": "Tried making the electrodes by ordering some films.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 5,
        "title": "\u7b2c5\u5929\uff1a\u4e0e\u6559\u6388\u548c\u5408\u4f5c\u8005\u5bf9\u63a5",
        "title_en": "Day 5: Syncing with Professor & Collaborators",
        "morning": "\u8ddf\u6559\u6388\u804a\u4e86\u804a\u8fd9\u51e0\u5929\u7684\u8fdb\u5c55\u3002\u91cd\u65b0\u60f3\u4e86\u60f3\u6bcf\u5468\u8ddf\u5408\u4f5c\u8005\u5f00\u4e00\u6b21\u4f1a\u6709\u6ca1\u6709\u5fc5\u8981\uff0c\u8fc7\u4e86\u4e00\u4e9b\u6280\u672f\u6307\u6807\u3002\u7ed9Multichannel Systems\u53d1\u4e86\u5c01\u90ae\u4ef6\uff0c\u60f3\u7ea6\u4e2a\u4f1a\u8bae\u8ba8\u8bba\u5b9a\u5236\u4fa7\u58c1\u65b9\u6848\u3002",
        "morning_en": "Talked to the professor about how everything went. Rethought whether meeting with the collaborators once per week is necessary, went over some specs. Sent an email to Multichannel Systems and asked for a meeting about a potentially custom-made side design.",
        "afternoon": "",
        "afternoon_en": "",
        "evening": "\u8ddf\u5b9e\u9a8c\u5ba4\u7684\u4eba\u51fa\u53bb\u559d\u5564\u9152\u3002\u5927\u5bb6\u804a\u4e86\u4e0d\u5c11\u8fd9\u51e0\u5929\u7684\u7cdf\u4e8b\uff0c\u7b11\u6b7b\u4e86\u3002",
        "evening_en": "Went out to drink beer with the gang. Had several laughs about their day.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 6,
        "title": "\u7b2c6\u5929\uff1aSolidWorks\u4e0e\u5893\u5730",
        "title_en": "Day 6: SolidWorks & the Cemetery",
        "morning": "\u7ee7\u7eed\u5b66SolidWorks\u3002\u65e0\u804a\u7684\u4e0a\u5348\u3002",
        "morning_en": "Continued learning SolidWorks. Boring morning.",
        "afternoon": "\u51fa\u53bb\u900f\u900f\u6c14\u3002\u5b66\u6821\u9644\u8fd1\u6709\u4e2a\u5893\u5730\uff0c\u5927\u5bb6\u4e4b\u524d\u63d0\u8fc7\uff0c\u53bb\u770b\u4e86\u770b\u3002",
        "afternoon_en": "Went out for some fresh air. There's a cemetery near the school\u2014people have mentioned it, so I went to check it out.",
        "evening": "\u5237\u4e86\u4e00\u5806styropyro\u7684\u89c6\u9891\u770b\u4e00\u4e9b\u79d1\u5e7b\u5411\u7684\u4e1c\u897f\uff0c\u5728Reddit\u4e0a\u7ffb\u5fae\u63a7\u5236\u5668\u7684\u6700\u65b0\u53d1\u660e\uff0c\u770b\u5230\u4e86\u4e00\u4e9b\u65b0\u82af\u7247\u3002\u5f00\u59cb\u7422\u78e8\u505a\u4e00\u4e2a\u81ea\u5b9a\u4e49\u7684\u98df\u7269\u70f9\u996a\u5668\u3002\u8bfb\u4e86\u51e0\u7bc7\u8bba\u6587\u3002",
        "evening_en": "Binge-watched styropyro online for some sci-fi stuff, scrolled on Reddit for the newest microcontroller inventions. Stumbled upon some new chips. Started thinking about making a custom food cooker. Read more papers.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 7,
        "title": "\u7b2c7\u5929\uff1a\u6362\u57f9\u517b\u6db2\u4e0e\u8ffd\u756a",
        "title_en": "Day 7: Sample Change & Anime Binge",
        "morning": "\u53bb\u5b9e\u9a8c\u5ba4\u5feb\u901f\u6362\u4e86\u4e00\u4e0b\u8111\u7ec4\u7ec7\u6837\u672c\u7684\u57f9\u517b\u6db2\u3002",
        "morning_en": "Went to the lab for a quick growth sample change for the brain sample.",
        "afternoon": "\u8eba\u5728\u5e8a\u4e0a\u770b\u4e86\u4e00\u6574\u5929\u52a8\u6f2b\u3002\u300a\u659a\uff01\u8d64\u7ea2\u4e4b\u7738\u300b(Akame ga Kill) 1-12\u96c6\u3002\u6211\u5f00\u59cb\u770b\u8fd9\u90e8\u662f\u56e0\u4e3a\u5973\u4e3b\u89d2\u957f\u5f97\u6709\u70b9\u50cf\u6211\uff0c\u800c\u4e14\u540d\u5b57\u91cc\u4e5f\u6709\"\u7ea2\"\u5b57\u3002",
        "afternoon_en": "Lay in bed all day watching anime. Akame ga Kill episodes 1-12. I started watching this because the protagonist girl kinda looks like me, and she has \"red\" in the name.",
        "evening": "\u4e3a\u591c\u88ad\u90e8\u961f\u7684\u5e0c\u5c14 (Sheele) \u611f\u5230\u7279\u522b\u96be\u8fc7\u3002\u4e00\u4e2a\u6f02\u6cca\u4e86\u4e00\u8f88\u5b50\u7684\u4eba\uff0c\u7ec8\u4e8e\u627e\u5230\u4e86\u4e00\u4e2a\u53ef\u4ee5\u79f0\u4e4b\u4e3a\u5bb6\u7684\u5730\u65b9\uff0c\u7136\u540e\u9a6c\u4e0a\u5c31\u6b7b\u4e86\u3002\u771f\u7684\u597d\u96be\u8fc7\u597d\u96be\u8fc7\u3002",
        "evening_en": "Felt really sad for Sheele. Someone who had drifted her whole life, finally found a place she calls home, and then immediately dies. I felt very very sad.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 8,
        "title": "\u7b2c8\u5929\uff1a\u8111\u7ec6\u80de\u57f9\u517b\u4e0e\u9a91\u884c\u793e",
        "title_en": "Day 8: Brain Cell Culture & Cycling Club",
        "morning": "\u8ddfMultichannel Systems\u5f00\u4e86\u7535\u8bdd\u4f1a\u8bae\u3002\u7136\u540e\u8ddf\u5b66\u957f\u53bb\u5b9e\u9a8c\u5ba4\uff0c\u7528\u8001\u914d\u65b9\u8bd5\u7740\u5728\u57f9\u517b\u76bf\u91cc\u57f9\u517b\u8111\u7ec6\u80de\u3002",
        "morning_en": "Had the meeting call with Multichannel Systems. Then went with senior to the lab, tried out the old recipe for growing brain cells in a petri dish.",
        "afternoon": "\u53bb\u4e0a\u8bfe\u3002\u5403\u4e86\u62c9\u9762\u3002",
        "afternoon_en": "Went to class. Ate ramen.",
        "evening": "\u7ee7\u7eed\u65e9\u4e0a\u7684\u8111\u7ec6\u80de\u57f9\u517b\uff0c\u62ff\u5230\u4e86\u751f\u7269\u5b9e\u9a8c\u5ba4\u5de5\u7a0b\u5e08\u7684lab write-off\u3002\u628a\u7b2c\u4e00\u7248\u7535\u6781\u710a\u76d8\u7684\u8bbe\u8ba1\u56fe\u53d1\u7ed9\u4e86\u6d01\u51c0\u5ba4\u5de5\u7a0b\u5e08\u53bb\u52a0\u5de5\u3002\n\n\u53c8\u53bb\u4e86\u9a91\u884c\u793e\uff0c\u53eb\u4e0a\u4e86Emily\u4e00\u8d77\u3002Emily\u8bb2\u4e86\u5979\u5b66\u9a91\u8f66\u7684\u7ecf\u5386\u3002",
        "evening_en": "Continued brain cell preparation from the morning, got the lab write-off from the bio-lab engineer. Emailed the design of the first pads to the cleanroom engineer for fabrication.\n\nWent to the cycling club again, called up Emily to go as well. Emily told me about her time learning to ride.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 9,
        "title": "\u7b2c9\u5929\uff1a\u5408\u4f5c\u8005\u4f1a\u8c08\u4e0e\u7ef4\u6301\u5927\u8111\u5b58\u6d3b",
        "title_en": "Day 9: Collaborator Meeting & Keeping the Brain Alive",
        "morning": "\u4eca\u5929\u8ddf\u5408\u4f5c\u8005\u5434\u535a\u58eb\u804a\u4e86\u804a\u3002\u4ed6\u8bf4\u4ed6\u4eec\u53ef\u4ee5\u6bd4\u8f83\u7a33\u5b9a\u5730\u671d\u4e00\u4e2a\u65b9\u5411\u53d1\u5c04\u4e00\u5927\u5806\u03bc\u5b50\u3002\u4f46\u662f\u03bc\u5b50\u8ddf\u7269\u8d28\u7684\u76f8\u4e92\u4f5c\u7528\u592a\u5f31\u4e86\uff0c\u6211\u4eec\u53ef\u80fd\u5f97\u628a\u5b9e\u9a8c\u642c\u5230\u7c92\u5b50\u52a0\u901f\u5668\u65c1\u8fb9\u3002\u4ed6\u4eec\u53ef\u4ee5\u8bd5\u7740\u628a\u8bbe\u5907\u642c\u5230\u6211\u4eec\u5b9e\u9a8c\u5ba4\uff0c\u4f46\u4f9b\u7535\u662f\u4e2a\u95ee\u9898\u3002\u800c\u4e14\u6211\u4eec\u5927\u6982\u8fd8\u8981\u4e24\u4e2a\u6708\u624d\u80fd\u5f00\u59cb\u91c7\u6570\u636e\u3002",
        "morning_en": "Today I talked with collaborator Dr. Wu. They said they could shoot a barrage of muons pretty reliably in one direction. But since muons are pretty non-interactive, we probably need to bring our experiment next to the particle accelerator. They can try bringing their equipment to our lab but the electricity will be a problem. We will also likely be 2 months out from gathering data.",
        "afternoon": "\u8bf4\u670d\u6559\u6388\u4ece\u57f9\u8bad\u4e2d\u5f04\u6765\u5927\u91cf\u6df7\u5408\u5e72\u7ec6\u80de\u3002\u5b9e\u9a8c\u5ba4\u51e0\u4e4e\u4ec0\u4e48\u90fd\u6709\uff0c\u6bd4\u5c0f\u5b9e\u9a8c\u5ba4\u597d\u591a\u4e86\u3002\u7ef4\u6301\u6837\u672c\u5b58\u6d3b\u7684\u8bbe\u5907\u90fd\u9f50\u5168\u3002\u7ee7\u7eed\u6628\u5929\u7684\u57f9\u517b\u51c6\u5907\u5de5\u4f5c\u3002\u56e0\u4e3a\u8981\u957f\u671f\u4f7f\u7528\u540c\u4e00\u4e2a\u6837\u672c\uff0c\u9700\u8981\u4e00\u4e2a\u81ea\u7ef4\u6301\u7684\u65b9\u6848\u6765\u8ba9\u5927\u8111\u6d3b\u7740\u3002\u6211\u4eec\u5728\u8003\u8651\u627e\u53e6\u4e00\u4e2a\u4e13\u95e8\u505a\u8111\u7ec7\u7ec7\u7ef4\u6301\u7684\u5408\u4f5c\u8005\u3002\n\n\u95ef\u8fdb\u4e01\u6559\u6388\u529e\u516c\u5ba4\u95ee\u4ed6\u8fd9\u4e8b\u3002\u4ed6\u6b63\u597d\u6709\u7a7a\uff0c\u592a\u597d\u4e86\u3002\u6211\u4eec\u8ba8\u8bba\u4e86\u5ef6\u957f\u81ea\u7ef4\u6301\u673a\u5236\u7684\u6d41\u7a0b\uff0c\u51b3\u5b9a\u5148\u4ece\u8fd9\u4e2a\u5165\u624b\u3002",
        "afternoon_en": "Convinced the professor to get a large amount of blended stem cells from training. The lab has almost everything I need\u2014much better than smaller labs. The equipment to keep things alive is all here. Continued yesterday's preparations for growing. Since we'll be spending a lot of time with one sample, we need a self-sustaining way to keep the brain alive. We're thinking of getting another collaborator specialized in keeping brains alive.\n\nBarged into Professor Ding's door to ask about this stuff. He had time\u2014wonderful. We talked over the procedure for the prolonged self-sustaining mechanism, and we decided to start working on that first.",
        "evening": "\u8bfb\u4e86\u4e24\u7bc7\u5173\u4e8e\u5982\u4f55\u7ef4\u6301\u8111\u7ec4\u7ec7\u5b58\u6d3b\u548c\u6a21\u62df\u751f\u957f\u6db2\u73af\u5883\u7684\u8bba\u6587\u3002\n\n\u770b\u4e86\u659a\u8d64\u7ea2\u4e4b\u779313-16\u96c6\u3002",
        "evening_en": "Read two papers on how to keep brain tissue alive and simulate the growth solution environment.\n\nWatched episodes 13-16 of Akame ga Kill.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 10,
        "title": "\u7b2c10\u5929\uff1a\u704c\u6d41\u6cf5\u8bbe\u8ba1",
        "title_en": "Day 10: Perfusion Pump Design",
        "morning": "\u6559\u6388\u8bf4\u4ed6\u8054\u7cfb\u4e86\u9694\u58c1\u505a\u4f53\u5916\u80da\u80ce\u7684\u540c\u4e8b\uff0c\u628a\u6211\u4ecb\u7ecd\u7ed9\u4e86\u5c0f\u8d3e\u3002\u5c0f\u8d3e\u8bf4\u5982\u679c\u6211\u60f3\u7528\u6df7\u5408\u8111\u7ec6\u80de\u7ec3\u4e60\u7684\u8bdd\uff0c\u4ed6\u660e\u5929\u65e9\u4e0a\u5c31\u53ef\u4ee5\u6765\u3002\u6240\u4ee5\u4eca\u5929\u6211\u5f97\u8bbe\u8ba1\u4e00\u4e2a\u53cc\u5411\u65e0\u83cc\u6cf5\u6765\u8f93\u9001\u57f9\u517b\u6db2\u3002\n\n\u57fa\u672c\u4e0a\u4e00\u6574\u5929\u90fd\u5728\u641e\u8fd9\u4e2a\u6cf5\u3002\u65b0\u7684\u5c0f\u6cf5\u8981\u4e00\u5468\u624d\u5230\uff0c\u4f46\u5c0f\u8d3e\u628a\u4ed6\u7684\u6cf5\u501f\u7ed9\u6211\u4e86\uff0c\u8fd9\u6837\u6211\u53ef\u4ee5\u7528\u7b2c4\u5929\u7684\u7ec6\u80de\u6765\u5b9e\u9a8c\u3002",
        "morning_en": "Professor said he got in contact with a colleague from the university right here doing in vitro embryos, and hooked me up with Jia. Jia said he could show up tomorrow morning if I wanted to practice with multi-blend brain cells, so today I have to design a 2-way sterilized pump for the growth fluid.\n\nPretty much the whole day I was working on the pump. We have to buy new small pumps that arrive in 1 week, but Jia gave me the pump from his setup so I can experiment with the cells from day 4.",
        "afternoon": "\u53bb\u4e0a\u8bfe\u3002\u4e0b\u8bfe\u540e\u8ddf\u4e00\u4e2a\u540c\u5b66\u804a\u4e86\u804a\uff0c\u6211\u4eec\u51b3\u5b9a\u4e00\u8d77\u5403\u5348\u996d\u3002\u4eba\u631a\u597d\u7684\u3002\u7ea6\u597d\u4e86\u5468\u4e94\u4e00\u8d77\u5199\u4f5c\u4e1a\u3002",
        "afternoon_en": "Went to class. After I finished class I talked to a classmate. We decided to get lunch together. He's pretty cool. We planned on finishing the homework together on Friday.",
        "evening": "\u62ff\u4e86\u4e2a\u73bb\u7483\u57f9\u517b\u76bf\uff0c\u628a\u4e24\u4e2a\u6cf5\u548c\u5206\u79bb\u819c\u90fd\u88c5\u4e0a\u53bb\u4e86\u3002\u770b\u8d77\u6765\u5f88\u5c71\u5be8\uff0c\u5e0c\u671b\u5c0f\u8d3e\u660e\u5929\u770b\u5230\u80fd\u63a5\u53d7🤞",
        "evening_en": "I got a glass petri dish, slapped two pumps on it as well as a separation membrane. It looks jank\u2014I hope Jia likes it tomorrow, fingers crossed.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 11,
        "title": "\u7b2c11\u5929\uff1akeeping_it_alive\u9879\u76ee\u542f\u52a8",
        "title_en": "Day 11: Project keeping_it_alive is a Go",
        "morning": "keeping_it_alive\u9879\u76ee\u6b63\u5f0f\u542f\u52a8\u3002\u867d\u7136\u5c3a\u5bf8\u6bd4\u5b9e\u9645\u768415\u00d71\u00d71\u5c0f\u5f88\u591a\uff0c\u4f46\u6cf5\u7406\u8bba\u4e0a\u5e94\u8be5\u591f\u7528\u3002\u5c0f\u8d3e\u770b\u5230\u6211\u7684\u88c5\u7f6e\u7b11\u4e86\u3002\u704c\u6d41\u8154\u914d\u6cf5\u548c\u6c14\u6ce1\u6355\u96c6\u5668\u8fd8\u6709\u5f88\u591a\u95e8\u9053\u3002\u5e78\u597d\u4ed6\u6709\u73b0\u6210\u7684\uff0c\u5efa\u8bae\u6211\u4eec\u5148\u81ea\u5df1\u628a\u5927\u6837\u672c\u957f\u51fa\u6765\u5c31\u884c\u3002\u771f\u662f\u4e2a\u9760\u8c31\u7684\u4eba\u3002\u6211\u592a\u559c\u6b22\u8fd9\u4e2a\u4eba\u4e86\u3002\n\n\u6d01\u51c0\u5ba4\u5de5\u7a0b\u5e08\u505a\u5b8c\u4e86\u7535\u6781\u58c1\uff0c\u6211\u628a\u88c5\u7f6e\u7ec4\u88c5\u597d\u4e86\u3002",
        "morning_en": "Project keeping_it_alive is a go. Although the dimensions are significantly smaller than the actual 15\u00d71\u00d71 size, the pump theoretically should be good enough to pump nutrients in. Jia saw the setup and chuckled. There were a lot more things associated with a perfusion chamber\u2014pumps and bubble traps. Luckily he has readily available ones, and recommended we just grow the large sample size on our own first. Quite a knowledgeable guy. I love this guy.\n\nThe cleanroom engineers finished the electrode walls and I assembled the setup.",
        "afternoon": "",
        "afternoon_en": "",
        "evening": "",
        "evening_en": "",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 12,
        "title": "\u7b2c12\u5929\uff1a\u6d4b\u8bd5\u4e0e\u5468\u603b\u7ed3",
        "title_en": "Day 12: Testing & Weekly Wrap-up",
        "morning": "\u53bb\u627e\u6d01\u51c0\u5ba4\u5de5\u7a0b\u5e08\uff0c\u62ff\u5230\u6837\u54c1\u3002\n\u7528\u663e\u5fae\u955c\u4e0b\u7684\u8d85\u7ec6\u5bfc\u7ebf\u6d4b\u8bd5\u4e86\u6bcf\u4e2a\u5f15\u811a\u7684\u63a5\u89e6\uff0c\u90fd\u6ca1\u95ee\u9898\u3002\u4ece\u5b9e\u9a8c\u5ba4\"\u501f\"\u4e86\u4e00\u6839\u6392\u7ebf\u3002\u9a8c\u8bc1\u82b1\u4e86\u4e0d\u5c11\u65f6\u95f4\u3002\n\n\u8ddf\u6559\u6388\u804a\u4e86\u9879\u76ee\u8fdb\u5c55\u3002\u8fd9\u5468\u4e3b\u8981\u8ddf\u5b66\u957f\u4e00\u8d77\u8bad\u7ec3\u3002\u7535\u6781\u7247\u7684\u52a0\u5de5\u505a\u5b8c\u4e86\uff0c\u4e5f\u6d4b\u8bd5\u4e86\u3002\u4e4b\u524d\u7684\u57f9\u517b\u6837\u672c\u5feb\u957f\u597d\u4e86\uff0c\u5f97\u53bb\u770b\u770b\u3002\u8fd9\u5468\u6700\u5927\u7684\u6210\u5c31\u5c31\u662f\u5230\u5904\"\u501f\"\u4e1c\u897f\u2014\u2014\u9694\u58c1\u5b66\u957f\u7684\u53cc\u6cf5\u704c\u6d41\u88c5\u7f6e\u3001\u6392\u7ebf\u3001\u53cc\u9762\u73bb\u7483\u80cc\u66b4\u9732\u63a5\u5934\u7684\u8fde\u63a5\u5668\u3002\u5e0c\u671b\u6240\u6709\u795e\u7ecf\u5143\u8fd9\u5468\u90fd\u80fd\u6d3b\u4e0b\u6765\u3002\u4e0d\u8fc7\u8fd8\u4e0d\u786e\u5b9a\u6211\u57f9\u517b\u7684\u795e\u7ecf\u5143\u591f\u4e0d\u591f\u957f\u671f\u4f7f\u7528\u3002\u6253\u7b97\u518d\u8ba2\u4e00\u4e9b\u661f\u5f62\u80f6\u8d28\u7ec6\u80de\u548c\u5c0f\u80f6\u8d28\u7ec6\u80de\u6765\u7a33\u5b9a\u60ac\u6d6e\u7cfb\u7edf\u3002\uff08\u6211\u4eec\u8ba1\u5212\u628a\u4e00\u90e8\u5206\u7528\u8584\u819c\u8d34\u5728\u58c1\u4e0a\uff0c\u5176\u4f59\u7684\u6254\u8fdb\u6eb6\u6db2\u91cc\u505a3D\u6620\u5c04\u7cfb\u7edf\u3002\u53ef\u80fd\u4f1a\u51fa\u73b0\u53ea\u6709\u8d34\u58c1\u7ec6\u80de\u624d\u6709\u8111\u529f\u80fd\u7684\u95ee\u9898\uff0c\u4f46\u5e94\u8be5\u6ca1\u5927\u788d\uff09",
        "morning_en": "Greeted the cleanroom engineer, got the sample.\nTested each pin's contacts with a super small wire viewed from the microscope. They all appeared well. Stole a ribbon cable from the lab. This took a while to verify.\n\nWent to the professor to talk about the project. This week I mainly trained with the senior. We got the fab done for the slides and tested them. Our previous samples are almost done growing\u2014I'll have to check up on them. This week's huge shoutout goes to stealing things from different labs\u2014like the duo pump diffusion setup from the senior next door, the ribbon cables, and the connectors for the back-exposed double glass setup. Ideally all the neurons survive this week, although I'm still not sure whether the neurons I grow will be effective enough for long term. I'm gonna go ahead and order more astrocytes and microglia for stabilizing a floating system. (We're planning on sticking some on the walls with films, then throwing the rest into the solution for a 3-D mapped system. We might have problems with brain functions only working from cells stuck to the cell walls, but it's probably fine.)",
        "afternoon": "\u628a\u4e00\u4e9b\u7ec6\u80de\u653e\u5230\u7535\u6781\u58c1\u7684\u4e00\u9762\u4e0a\u3002\u660e\u5929\u770b\u6548\u679c\u3002",
        "afternoon_en": "Placed some cells on one side of one cell wall electrode. We'll see how it works tomorrow.",
        "evening": "\u53c8\u8ddf\u6559\u6388\u53bb\u559d\u9152\u4e86\u3002\u5c0f\u8d3e\u8bf4\u4e86\u4e2a\u6709\u8da3\u7684\u89c2\u70b9\u2014\u2014\u5f00\u5b66\u524d\u4e24\u5468\u662f\u8ba4\u8bc6\u4eba\u7684\u6700\u4f73\u65f6\u673a\uff0c\u4e24\u5468\u4e4b\u540e\u5927\u5bb6\u5c31\u6ca1\u90a3\u4e2a\u52a8\u529b\u53bb\u8ba4\u8bc6\u65b0\u4eba\u4e86\uff0c\u66f4\u522b\u8bf4\u4e3b\u52a8\u793e\u4ea4\u4e86\u3002\u5475\uff0c\u8fd9\u54e5\u4eec\u4f30\u8ba1\u5c31\u662f\u5728\u5230\u5904\u642d\u8baa\u7f8e\u5973\u5427\u563f\u563f\u3002\u5c0f\u6591\u4eca\u5929\u631a\u597d\u7684\uff0c\u5e26\u4e86\u66f2\u5947\u997c\u5e72\uff0c\u8fd8\u8ddf\u6211\u4eec\u804a\u4e86\u5979\u5bf9Cosplay\u5236\u4f5c\u7684\u5174\u8da3\u3002\u7f1d\u7eab\u548c\u6570\u5b66\u95ee\u9898\u8fd8\u771f\u4e0d\u592a\u4e00\u6837\u3002\n\n\u6211\u4eec\u4eca\u5929\u804a\u4e86\u597d\u591a\u6709\u7684\u6ca1\u7684\u3002\u6bd4\u5982AI\u8981\u7edf\u6cbb\u4e16\u754c\u4ec0\u4e48\u7684\u3002\u8bf4\u6765\u4e5f\u5947\u602a\uff0c\u6211\u89c9\u5f97AI\u53d1\u5c55\u5f97\u53c8\u5feb\u53c8\u6162\u3002",
        "evening_en": "Went drinking with the professor again. Jia mentioned a fun thing\u2014how the first two weeks of school is the best time to meet people, as people after the first two weeks no longer have the drive to meet new people, much less meet people in general. Meh, the guy's probably just running around with a bunch of women hehe. Bian was pretty nice today. She brought cookies and told us about her interest in cosplay making. Sewing and math problems don't really have that much in common.\n\nWe talked over a lot of stupid shit today. Like how AI is gonna take over the world or whatever. It's weird\u2014I feel like it's evolving fast and slow at the same time.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 13,
        "title": "\u7b2c13\u5929\uff1a\u5c16\u5cf0\u4fe1\u53f7\u4e0ePong",
        "title_en": "Day 13: Spikes & Pong",
        "morning": "\u53bb\u68c0\u67e5\u4e86\u7b2c4\u5929\u7684\u7ec6\u80de\u57f9\u517b\u6837\u672c\uff0c\u8fd8\u6d3b\u7740\uff01\u592a\u597d\u4e86\u3002\u628a\u914d\u65b9\u5b58\u597d\uff0c\u51c6\u5907\u505a\u65b0\u7684\u57f9\u517b\u3002\u628a\u7ec6\u80de\u8f6c\u79fb\u5230\u4e86\u5355\u5e73\u9762\u7535\u6781\u4e0a\uff0c\u8fde\u63a5\u5230\u5b9e\u9a8c\u5ba4\u90a3\u53f0\u4e00\u4e07\u7f8e\u5143\u7684ADC\uff0c\u5c45\u7136\u770b\u5230\u4e86\u4e00\u4e9b\u5c16\u5cf0\u4fe1\u53f7\u3002\u8fdb\u5c55\u4e0d\u9519\u3002\n\n\u7b2c\u4e00\u4e2a\u60f3\u6cd5\u5c31\u662f\u505a\u4e00\u4e2a\u6253\u4e52\u4e53\u7403\u7684\u673a\u5668\uff08Pong\uff09\u3002\n\n\u6240\u4ee5\u4eca\u5929\u6211\u8981\u5199\u4e00\u4e2aPong\u6e38\u620f\u7684\u4ee3\u7801\uff0c\u5e0c\u671b\u80fd\u590d\u73b0\u90a3\u4e2a\u5b9e\u9a8c\u3002",
        "morning_en": "Went to check on the cell culture samples from day 4\u2014they're still alive. Great. We can store the recipe and move on to making new cultures. I transferred the cells to the single plane electrode and connected the system to our lab's $10k ADC, and I was able to see some spikes. Pretty good progress.\n\nFirst idea: build a Pong machine.\n\nSo today I'm gonna code Pong, and just hope we can replicate the setup.",
        "afternoon": "",
        "afternoon_en": "",
        "evening": "",
        "evening_en": "",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 14,
        "title": "\u7b2c14\u5929\uff1a\u795e\u7ecf\u5143\u4f1a\u6253Pong\u4e86\uff01",
        "title_en": "Day 14: Neurons Play Pong!",
        "morning": "\u6628\u5929Cursor\u57fa\u672c\u5e2e\u6211\u628a\u6574\u4e2aPong\u5199\u5b8c\u4e86\u3002\n\n\u5b9e\u9a8c\u8bbe\u8ba1\uff1a\n\u7b2c1\u7ec4\u7535\u6781\uff1a\n\u4e0a\u65b9\u7535\u6781 \u2192 \u7403\u5728\u9ad8\u5904\n\u4e2d\u95f4\u7535\u6781 \u2192 \u7403\u5728\u4e2d\u95f4\n\u4e0b\u65b9\u7535\u6781 \u2192 \u7403\u5728\u4f4e\u5904\n\n\u7b2c2\u7ec4\u7535\u6781\uff1a\n\u5de6\u4fa7\u7535\u6781\u653e\u7535\u591a \u2192 \u6321\u677f\u4e0a\u79fb\n\u53f3\u4fa7\u7535\u6781\u653e\u7535\u591a \u2192 \u6321\u677f\u4e0b\u79fb\n\n\u6321\u677f\u63a5\u4f4f\u7403\uff08\u597d\u7684\uff09\u2192 \u89c4\u5f8b\u6027\u8282\u5f8b\u8109\u51b2\n\u6321\u677f\u6ca1\u63a5\u4f4f\uff08\u574f\u7684\uff09\u2192 \u591a\u4e2a\u7535\u6781\u4e0a\u7684\u4e0d\u89c4\u5219\u7206\u53d1\n\n\u505a\u597d\u4e4b\u540e\u5c31\u628a\u5b83\u63a5\u5230\u57f9\u517b\u7269\u4e0a\u6a21\u62df\u4e86\u3002\uff08\u5230\u76ee\u524d\u4e3a\u6b62\u6211\u4eec\u8fdb\u5c55\u98de\u5feb\u3002\u6211\u4eec\u5047\u8bbe\u53ef\u4ee5\u4ece\u5b9e\u9a8c\u5ba4\u548c\u5bf9\u9762\u5b9e\u9a8c\u5ba4\"\u501f\"\u5230\u6240\u6709\u4e1c\u897f\u3001\u4ec0\u4e48\u90fd\u6709\u73b0\u6210\u7684\u3001\u5bfc\u5e08\u4eec\u90fd\u6709\u7a7a\u3001\u57f9\u517b\u7269\u7b2c\u4e00\u6b21\u5c31\u6ca1\u6b7b\u3001\u751f\u957f\u56e0\u5b50\u968f\u65f6\u53ef\u7528\u3001\u5b9e\u9a8c\u5ba4\u7a7a\u95f4\u4e5f\u4e0d\u53d7\u9650\u3002\u7b80\u76f4\u662f\u4ed9\u5883\u3002\u54e6\uff0c\u8fd8\u6709\u6d01\u51c0\u5ba4\u5de5\u7a0b\u5e08\u80fd\u57283\u5929\u5185\u5b8c\u6210\u5de5\u5355\u3002\uff09\n\n\u800c\u4e14\uff01\u771f\u7684\u6210\u529f\u4e86\uff01\u592a\u597d\u4e86\uff01",
        "morning_en": "Cursor practically coded the entire Pong yesterday.\n\nExperimental setup:\nGroup 1 electrodes:\nTop electrodes \u2192 ball high\nMiddle electrodes \u2192 ball center\nBottom electrodes \u2192 ball low\n\nGroup 2 electrodes:\nMore spikes left electrodes \u2192 move paddle up\nMore spikes right electrodes \u2192 move paddle down\n\nIf the paddle hits the ball (good) \u2192 regular rhythmic pulses\nIf the paddle misses (bad) \u2192 irregular bursts across many electrodes\n\nOk, so now we made this, I connected it to the culture to simulate. (So far we have been moving extremely quickly. We assume we can steal everything from our lab and the lab across, with everything available, the mentors always have free time, the cultures don't die on first try, the growth agents are readily available, and no lab space constraints. Wonderland over here. Oh, and the cleanroom engineer can finish something in 3 days' work order.)\n\nAnd it just works! Wonderful.",
        "afternoon": "",
        "afternoon_en": "",
        "evening": "\u7ed9\u6559\u6388\u53d1\u90ae\u4ef6\u544a\u8bc9\u4ed6\u8fd9\u4e2a\u6d88\u606f\uff0c\u660e\u5929\u8ba8\u8bba\u4e0b\u4e00\u6b65\u8ba1\u5212\u3002\n\n\u8fd8\u9700\u8981\u51c6\u5907\uff1a\u4e0b\u4e00\u6b65\u76843D\u7ec4\u88c5\u7cfb\u7edf\u3001\u505a\u6cd5\u62c9\u7b2c\u7b3c\u3001\u4e70\u66f4\u9ad8\u7ea7\u7684\u795e\u7ecf\u5143\uff08\u591a\u5df4\u80fa\u795e\u7ecf\u5143\u3001\u53d7\u4f53\u3001\u76ae\u5c42\u795e\u7ecf\u5143\u3001\u611f\u89c9\u795e\u7ecf\u5143\u7b49\uff09\u3002\u8981\u60f3\u4e2a\u65b9\u6848\u6765\u63d0\u5347\u7cfb\u7edf\u6027\u80fd\uff0c\u53ef\u80fd\u9700\u8981\u4e70\u66f4\u591a\u7684\u6a21\u62df\u8bfb\u51fa\u8bbe\u5907\uff0c\u8fd8\u6709\u4e0b\u5468\u53ef\u80fd\u51fa\u73b0\u7684\u5404\u79cd\u610f\u5916\u60c5\u51b5\u3002\n\n\u5f55\u4e86\u795e\u7ecf\u5143\u6253Pong\u7684\u89c6\u9891\u3002",
        "evening_en": "Emailed the professor to tell him about it so we can discuss further actions tomorrow.\n\nAlso need to prepare: next steps for the 3D assembly system, making the Faraday cage, buying more advanced neurons (dopamine neurons, receptors, cortical neurons, sensory neurons, etc.). Need to think of a scheme to improve system capabilities\u2014might have to buy a lot more analog readouts, and plan for several unforeseeable things for next week.\n\nRecorded the video of the neuron playing Pong.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 15,
        "title": "\u7b2c15\u5929\uff1a\u5411\u6559\u6388\u6c47\u62a5",
        "title_en": "Day 15: Showing the Professor",
        "morning": "\u65e9\u4e0a\u53bb\u627e\u4e01\u6559\u6388\u3002\u4ed6\u5fc3\u60c5\u4e0d\u9519\uff0c\u6211\u4e5f\u5fc3\u60c5\u4e0d\u9519\u3002\"\u563f\u6559\u6388\uff0c\u770b\u770b\u8fd9\u4e2a\u3002\"\n\n\u8bf4\u5b9e\u8bdd\u6211\u771f\u7684\u631a\u5f00\u5fc3\u7684\u3002\u4ed6\u770b\u8d77\u6765\u4e5f\u5f88\u60ca\u8bb6\u3002\u4e0d\u8fc7\u8fd9\u4e2a\u5b9e\u9a8c\u4ed6\u4ee5\u524d\u4e5f\u505a\u8fc7\uff08\u6240\u4ee5\u5b9e\u9a8c\u5ba4\u91cc\u624d\u6709\u90a3\u4e48\u591a\u73b0\u6210\u7684\u96f6\u4ef6\uff09\u3002\u4ed6\u95ee\u6211\u63a5\u4e0b\u6765\u6709\u4ec0\u4e48\u8ba1\u5212\u3002\u901a\u5e38\u8fd9\u4e9b\u5b9e\u9a8c\u4e0d\u4f1a\u8fdb\u5c55\u8fd9\u4e48\u987a\u5229\uff0c\u5c24\u5176\u662f\u5bf9\u4e00\u5e74\u7ea7\u535a\u58eb\u751f\u6765\u8bf4\u3002\u6211\u4eec\u8ba8\u8bba\u4e86\u5927\u8111\u60ac\u6d6e\u65b9\u6848\u2014\u2014\u662f\u7528\u6a21\u5f0f\u54cd\u5e94\u8fd8\u662f\u591a\u5df4\u80fa\u54cd\u5e94\u7cfb\u7edf\u3002\u603b\u4e4b\u8fd9\u6b21\u5c55\u793a\u8ba9\u4ed6\u5f88\u9ad8\u5174\uff0c\u800c\u4e14\u4e5f\u6ca1\u82b1\u4ed6\u592a\u591a\u94b1\uff0c\u6240\u4ee5\u7686\u5927\u6b22\u559c\u3002",
        "morning_en": "In the morning I went to talk with Professor Ding. He was in a pretty good mood. I was also in a pretty good mood. \"Hey prof, check this.\"\n\nI gotta tell you, I was pretty happy. He looked very impressed too. Unfortunately this was also something he had done before. (Hence we had all the components just lying around the lab.) He asked me what my plans were for the foreseeable future. Normally these experiments don't run so smoothly, especially among first-year PhD students. We discussed the suspension of the brain\u2014whether to use patterned response or dopamine response systems. But the show-off made him really happy and didn't cost him that much, so he's happy.",
        "afternoon": "\u53bb\u4e0a\u8bfe\u3002",
        "afternoon_en": "Went to class.",
        "evening": "\u56de\u5bb6\u4e86\u3002\u4eca\u5929\u631a\u7d2f\u7684\u3002\u53bb\u52a8\u6f2b\u793e\u770b\u770b\u6700\u8fd1\u5728\u641e\u4ec0\u4e48\u3002\u4eba\u6bd4\u4e4b\u524d\u5c11\u4e86\u597d\u591a\u3002\u6211\u731c\u4e3b\u8981\u539f\u56e0\u662f\u5927\u591a\u6570\u65b0\u4eba\u5f88\u5feb\u5c31\u4e0d\u6765\u4e86\uff0c\u6240\u4ee5\u8001\u6210\u5458\u4e5f\u4e0d\u600e\u4e48\u70ed\u60c5\u5730\u6b22\u8fce\u65b0\u4eba\u4e86\u3002\n\n\u53c8\u770b\u4e863\u96c6\u659a\u8d64\u7ea2\u4e4b\u779b\u3002\u604b\u7231\u573a\u9762\u631a\u4e0a\u5934\u7684\uff0c\u6253\u6597\u4e5f\u5f88\u7cbe\u5f69\u3002",
        "evening_en": "Went home. The day was pretty exhausting. Went to the anime club to see how things are going. They had a lot fewer people than before. I guess a key reason nobody treats newcomers with much welcome is because most newcomers just leave very often.\n\nWatched 3 more episodes of Akame ga Kill. The romance scene was pretty hot, and the fighting was fun.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 16,
        "title": "\u7b2c16\u5929\uff1a\u8ba2\u8d2d\u9ad8\u7ea7\u7ec6\u80de",
        "title_en": "Day 16: Ordering Advanced Cells",
        "morning": "\u65e2\u7136\u0032D\u5b9e\u9a8c\u6210\u529f\u4e86\uff0c\u4eca\u5929\u8ba2\u4e86\u0038\u79cd\u4e0d\u540c\u7684\u7ec6\u80de\u7c7b\u578b\uff0c\u82b1\u4e86\u6559\u6388\u0031\u0030\u4e07\u7f8e\u5143\u3002\u4ed6\u7729\u7740\u773c\u770b\u6211\uff0c\u4f46\u540e\u6765\u4ed6\u95ee\u4e86\u95ee\u9694\u58c1\u5b9e\u9a8c\u5ba4\uff0c\u4ed6\u4eec\u8bf4\u4ee5\u540e\u53ef\u80fd\u4e5f\u7528\u5f97\u4e0a\uff0c\u6240\u4ee5\u53ef\u4ee5\u5408\u4e70\u3002\u6211\u5c31\u662f\u5e86\u5e78\u8fd9\u4e8b\u6210\u4e86\u3002\u4eca\u5929\u6ca1\u8ddf\u5408\u4f5c\u8005\u901a\u8bdd\uff0c\u7ed9\u4ed6\u53d1\u4e86\u90ae\u4ef6\u4ecb\u7ecd\u6211\u4eec\u7684\u5b9e\u9a8c\u88c5\u7f6e\uff0c\u4ed6\u606d\u559c\u4e86\u6211\u3002\n\n\u65b0\u8bbe\u8ba1\u53d7\u5c0f\u8d3e\u7cfb\u7edf\u7684\u542f\u53d1\u2014\u2014\u4ece\u4e24\u4fa7\u63a5\u5165\u704c\u6d41\u6cf5\uff0c\u6cf5\u4ece\u5404\u6bb5\u4e2d\u95f4\u4f9b\u7ed9\u3002\u4e4b\u524d\u7684\u8bbe\u8ba1\u6ca1\u8003\u8651\u6cf5\u7684\u4f4d\u7f6e\uff0c\u6240\u4ee5\u8981\u8ba9\u6d01\u51c0\u5ba4\u5de5\u7a0b\u5e08\u91cd\u505a\u4e00\u7248\u3002\u0032\u0063\u006d\u00d7\u0032\u0063\u006d\u4e5f\u592a\u5c0f\u4e86\uff0c\u6253\u7b97\u6269\u5230\u0034\u00d7\u0034\uff0c\u56e0\u4e3a\u6cf5\u9700\u8981\u6269\u6563\u7a7a\u95f4\uff0c\u0031\u0035\u003a\u0032\u7684\u6bd4\u4f8b\u771f\u7684\u4f1a\u5f71\u54cd\u7cfb\u7edf\u6548\u679c\u3002\n\n\u6240\u4ee5\u4eca\u5929\u5f00\u59cb\u505a\u65b0\u7248\u672c\u0031\u0035\u003a\u0032\u003a\u0032\u7684\u8bbe\u8ba1\uff0c\u4e70\u7ec6\u80de\u3001\u65b0\u6cf5\u3001\u6c14\u6ce1\u6355\u96c6\u5668\u3001\u7ba1\u8def\u3001\u65b0\u771f\u7a7a\u6cf5\u3001\u65b0ADC\u2026\u2026\u4e00\u5806\u4e1c\u897f\u3002\u4e01\u6559\u6388\u770b\u5230\u6211\u7ed9\u4ed6\u7684\u6e05\u5355\u60ca\u5446\u4e86\uff08\u6211\u4e5f\u82b1\u4e86\u4e0d\u5c11\u529f\u592b\u5217\u8fd9\u4e2a\u6e05\u5355\u7684\uff0c\u6211\u4e5f\u4e0d\u60f3\u82b1\u94b1\u554a\uff09",
        "morning_en": "Since the 2D experiment worked, today I ordered 8 different cell types that cost our professor $100k. He had this squinted eyes look on his face, but then he asked around and neighboring labs said they may need it later as well, so it would be a pooled buy. I'm just glad it worked. I didn't talk with the collaborator today\u2014sent him an email of our setup, and he congratulated me on the task.\n\nMy new design draws inspiration from Jia's system\u2014we take in the diffusion pump from both sides, and host the pump from the middle of the sections. The previous design did not account for pump locations, so we'll have to get the cleanroom engineers to make a new version. The 2cm\u00d72cm is also too small\u2014I'll expand it to 4\u00d74 because the pumps need to diffuse, and having a 15:2 ratio really undermines the effects I imagine our system will hold.\n\nSo today I'll get started with the new version 15:2:2, buy cells, buy a new pump, bubble stopper, tubes, a new vacuum pump, a new ADC, and yeah. Professor Ding was shocked looking at the long list of items I gave him (I worked kinda hard on it too\u2014I wish I didn't have to buy either).",
        "afternoon": "",
        "afternoon_en": "",
        "evening": "",
        "evening_en": "",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 17,
        "title": "\u7b2c17\u5929\uff1a\u8bbe\u8ba1\u65e5",
        "title_en": "Day 17: Design Day",
        "morning": "\u4eca\u5929\u4e3b\u8981\u5c31\u662f\u505a\u8bbe\u8ba1\u3002\u6837\u672c\u8fd8\u8981\u4e00\u4e2a\u6708\u624d\u5230\uff0c\u6240\u4ee5\u5728\u8fd9\u671f\u95f4\u5c31\u8bfb\u8bfb\u8bba\u6587\u5427\u3002\u51b5\u4e14\u641e3D\u7cfb\u7edf\u4e5f\u4e0d\u5bb9\u6613\u3002",
        "morning_en": "Today I mainly worked on the designs. The samples aren't gonna come in for another month. So in the meantime I'll just read some papers. Plus, growing a 3D system wouldn't be easy.",
        "afternoon": "\u4e0a\u8bfe\u3002",
        "afternoon_en": "Class.",
        "evening": "\u628a\u4e4b\u524d\u768415\u00d72\u00d72\u7ec4\u88c5\u8d77\u6765\u4e86\u3002\u770b\u8d77\u6765\u8fd8\u884c\u3002\u6253\u7b97\u5728\u91cc\u9762\u8bd5\u7740\u957f\u4e00\u4e9b\u7ec6\u80de\uff0c\u4f46\u4e00\u6279\u7684\u6210\u672c\u4e0d\u4f4e\uff0c\u5f97\u5c0f\u5fc3\u70b9\u3002",
        "evening_en": "Put together the previous 15\u00d72\u00d72. It's holding up well. I'll try growing some cells in here, but it costs a lot of money for a batch, so I guess I'll have to be careful.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 18,
        "title": "\u7b2c18\u5929\uff1a\u914d\u6eb6\u6db2",
        "title_en": "Day 18: Solution Prep",
        "morning": "\u4e3b\u8981\u5728\u914d\u6eb6\u6db2\uff0c\u51c6\u5907\u65b0\u4e00\u6279\u57f9\u517b\u3002\u4e4b\u524d\u4ece\u5c0f\u8d3e\u90a3\u501f\u7684\u88c5\u7f6e\u4ed6\u8981\u7528\uff0c\u53ea\u597d\u53c8\u53bb\u501f\u4e86\u4e00\u6b21\u3002\u8fd9\u88c5\u7f6e\u770b\u8d77\u6765\u592a\u5c71\u5be8\u4e86\u2014\u2014\u5c31\u662f\u4e00\u6761\u6c34\u69fd\u4e24\u5934\u5404\u4e00\u4e2a\u6cf5\uff0c\u770b\u8d77\u6765\u597d\u641e\u7b11\u3002\u628a\u6392\u7ebf\u63a5\u5230\u7535\u6781\u8ff9\u7ebf\u4e0a\uff0c\u8fde\u4e0a\u7535\u8111\uff0c\u4fe1\u53f7\u770b\u8d77\u6765\u6ca1\u95ee\u9898\u3002\u597d\u4e86\uff0c\u57f9\u517b\u7269\u653e\u8fdb\u53bb\u4e86\uff0c\u63a5\u4e0b\u6765\u5c31\u7b49\u7740\u5427\u3002",
        "morning_en": "Mainly prepared the solutions, the preparations for growing a new batch. I had kept the previous setup from Jia but he had to use it earlier, so I had to go borrow from him again. The setup looks so jank. It's a slit of water with two pumps at the ends. It looks so wacky. I slapped the ribbon cables to the electrical traces, connected it to the computer, and the signals looked fine. Well, the cultures are in the slit now. All I have to do is wait.",
        "afternoon": "",
        "afternoon_en": "",
        "evening": "",
        "evening_en": "",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 19,
        "title": "\u7b2c19\u5929\uff1a\u8bba\u6587\u65e5\u4e0e\u8e22\u8e0f\u821e",
        "title_en": "Day 19: Paper Day & Tap Dancing",
        "morning": "\u8bfb\u8bba\u6587\u65e5\u3002",
        "morning_en": "Paper day.",
        "afternoon": "",
        "afternoon_en": "",
        "evening": "\u665a\u4e0a\u53bb\u9152\u5427\u8fd8\u662f\u5f88\u5f00\u5fc3\u3002\u6211\u4eec\u5f00\u59cb\u8dd1\u9898\u804a\u6559\u6388\u7684\u611f\u60c5\u751f\u6d3b\u4e86\u3002\u4ed6\u592a\u5e78\u798f\u4e86\u53cd\u800c\u6709\u70b9\u65e0\u804a\u3002\u5c0f\u8d3e\u53c8\u5728\u90a3\u5927\u79c0\u7279\u79c0\uff0c\u8bf4\u4ed6\u73b0\u5728\u5f97\u8003\u8651\u600e\u4e48\u5c55\u73b0\u81ea\u5df1\u9ad8\u5927\u4e0a\u7684\u4e00\u9762\uff0c\u8fd8\u6709\u4ed6\u7684\u5b9e\u9a8c\u591a\u4e48exciting\uff08\u6211\u4e0d\u7406\u89e3\u4ed6\u4e3a\u4ec0\u4e48\u8fd9\u4e48\u5174\u594b\uff0c\u4ed6\u7684\u8fdb\u5ea6\u6bd4\u6211\u6162\u4e24\u767e\u4e07\u5e74\uff0c\u6240\u4ee5\u6211\u624d\u8001\u662f\u7528\u4ed6\u7684\u751f\u7269\u88c5\u7f6e\uff09\u3002\u5475\u3002\n\n\u4e01\u6559\u6388\u4f1a\u8e22\u8e0f\u821e\uff01\u4eca\u5929\u770b\u4ed6\u8df3\u8e22\u8e0f\u821e\u771f\u7684\u592a\u9177\u4e86\u3002\u6211\u5f97\u627e\u4e2a\u65f6\u95f4\u8ba9\u4ed6\u6559\u6559\u6211\u3002",
        "evening_en": "In the evening the bar was still great. We started going off topic into the professor's love life. It was so happy it was kinda boring. Jia put on this huge show of how he now has to decide how to portray himself as high and mighty, as well as how his experiment is so exciting (which I don't understand why he's so excited, as he is practically 2 million years slower than me\u2014hence why I use his bio setup so often). Meh.\n\nProfessor Ding can tap dance! It was really cool to see him tap dance today. Man, I gotta get him to teach me sometime.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
    {
        "day": 20,
        "title": "\u7b2c20\u5929\uff1a\u8bba\u6587\u65e5\u4e0e\u611f\u609f",
        "title_en": "Day 20: Paper Day & Reflections",
        "morning": "\u8bfb\u8bba\u6587\u65e5#2\u3002",
        "morning_en": "Paper day #2.",
        "afternoon": "",
        "afternoon_en": "",
        "evening": "\u564e\uff0c\u4eca\u5929\u770b\u5b8c\u4e86\u659a\u8d64\u7ea2\u4e4b\u779b\u3002\u8fd9\u4e9b\u4eba\u90fd\u6709\u81ea\u5df1\u5728\u8fd9\u4e2a\u4e16\u754c\u4e0a\u7684\u4fe1\u5ff5\uff0c\u4e00\u76f4\u8d70\u5230\u6700\u540e\u3002\u6211\u559c\u6b22\u8fd9\u6837\u7684\u4eba\u3002\u6709\u65f6\u5019\u6211\u5e0c\u671b\u81ea\u5df1\u4e5f\u80fd\u8fd9\u6837\u3002\u8d70\u5728\u4e00\u7fa4\u6709\u7406\u60f3\u7684\u4eba\u4e2d\u95f4\uff0c\u6070\u597d\u8d70\u7740\u540c\u4e00\u6761\u8def\u3002\n\n\u5b9e\u9a8c\u5ba4\u7684\u4eba\u90fd\u4e0d\u9519\uff0c\u4f46\u603b\u89c9\u5f97\u6709\u8ddd\u79bb\u3002\u597d\u50cf\u4e0d\u80fd\u5206\u4eab\u4e00\u4e9b\u65e0\u5398\u5934\u7684\u4e1c\u897f\u3002\u4e5f\u6ca1\u529e\u6cd5\u8ddf\u522b\u4eba\u804a\u8fd9\u4e9b\u3002\u4e5f\u8bb8\u7f51\u4e0a\u80fd\u627e\u5230\u4eba\u5427\u3002",
        "evening_en": "Oh, I finished Akame ga Kill today. Wow, these guys all had their own individual rights in the world, and walked their paths till the very end. I like people who do this. I wish sometimes I can be like this. Walking among people who have an ideal they chase and just happen to walk the same path we do.\n\nPeople at the lab are cool, but they also feel distant. Like I can't share stupid things. I can't talk to other people about this. Maybe someone online will do.",
        "night": "",
        "night_en": "",
        "noon": "",
        "noon_en": ""
    },
]

# Keep existing days 1, 2, 35, 42, 51 from the original
existing = data['kurisu']['schedules']
keep_days = {1, 2, 35, 42, 51}
kept = [s for s in existing if s['day'] in keep_days]

# Combine: kept originals + new days 3-20
all_schedules = kept + new_days

# Sort by day number
all_schedules.sort(key=lambda x: x['day'])

data['kurisu']['schedules'] = all_schedules

# ============================================================
# UPDATE SCIENTIFIC ANCHORS (Months 1 and 2) for consistency
# ============================================================

anchors = data['kurisu']['story']['scientificAnchors']

for a in anchors:
    if a['month'] == 1:
        a['title'] = "\u7b2c1\u6708\uff1a\u4ece\u96f6\u5230Pong"
        a['title_en'] = "Month 1: From Zero to Pong"
        a['content'] = (
            "\u8fdb\u5165\u5b9e\u9a8c\u5ba4\uff0c\u7ed3\u8bc6\u65b0\u540c\u4e8b\uff08\u5c0f\u8d3e\u3001\u5c0f\u4f0a\u3001\u5c0f\u6591\uff09\uff0c\u52a0\u5165\u9a91\u884c\u793e\u548c\u52a8\u6f2b\u793e\u3002\n"
            "\u8bbe\u8ba1\u4e86\u4e24\u9636\u6bb5\u5b9e\u9a8c\u7cfb\u7edf\uff1a\n"
            "\u7b2c\u4e00\u9636\u6bb5\u2014\u2014\"\u7535\u6781\u6d4b\u8bd5\u5957\u4ef6\"\uff1a\u6cd5\u62c9\u7b2c\u7b3c + \u8111\u7ec6\u80de\u57f9\u517b + \u4fa7\u58c1\u7535\u6781 + \u4fe1\u53f7\u8bfb\u51fa\uff0c\u901a\u8fc7\u8ba9\u57f9\u517b\u7684\u5927\u8111\u73a9Doom\u6765\u9a8c\u8bc1\u529f\u80fd\u3002\n"
            "\u7b2c\u4e8c\u9636\u6bb5\u2014\u2014\u91cf\u5b50\u4fe1\u53f7\u9a8c\u8bc1\u5668\uff1a\u8ba9\u9694\u58c1\u52a0\u901f\u5668\u5b9e\u9a8c\u5ba4\u53d1\u5c04\u03bc\u5b50\u675f\uff0c\u7528\u5207\u4f26\u79d1\u592b\u63a2\u6d4b\u5668\u8bfb\u4e00\u7ef4\u4fe1\u53f7\uff0c\u9a8c\u8bc1\u8111\u7ec4\u7ec7\u5728\u91cf\u5b50\u4fe1\u53f7\u4e0b\u662f\u5426\u4ecd\u80fd\u6b63\u5e38\u8fd0\u4f5c\u3002\n\n"
            "\u5b66\u4e86SolidWorks\u5efa\u6a21\uff08\u53d1\u73b0\u9700\u8981\u751f\u7269\u517c\u5bb9\u6750\u6599\uff09\uff0c\u5b8c\u6210\u751f\u7269\u5b9e\u9a8c\u5ba4\u5b89\u5168\u57f9\u8bad\u3002\u7528\u8001\u914d\u65b9\u5728\u57f9\u517b\u76bf\u91cc\u57f9\u517b\u4e86\u7b2c\u4e00\u6279\u8111\u7ec6\u80de\uff0c\u6d01\u51c0\u5ba4\u5de5\u7a0b\u5e083\u5929\u5185\u52a0\u5de5\u5b8c\u6210\u7535\u6781\u58c1\u3002\n"
            "\u4e0e\u5408\u4f5c\u8005\u5434\u535a\u58eb\u521d\u6b65\u8ba8\u8bba\u4e86\u03bc\u5b50\u675f\u7684\u63a5\u5165\u65b9\u6848\uff1a\u4ed6\u4eec\u53ef\u4ee5\u7a33\u5b9a\u53d1\u5c04\u03bc\u5b50\u675f\uff0c\u4f46\u7531\u4e8e\u03bc\u5b50\u4e0e\u7269\u8d28\u76f8\u4e92\u4f5c\u7528\u5f31\uff0c\u53ef\u80fd\u9700\u8981\u5c06\u5b9e\u9a8c\u642c\u5230\u52a0\u901f\u5668\u65c1\u8fb9\uff0c\u9884\u8ba1\u4e24\u4e2a\u6708\u540e\u624d\u80fd\u5f00\u59cb\u91c7\u6570\u636e\u3002\n\n"
            "\u91cd\u5927\u7a81\u7834\uff1a\u6210\u529f\u8ba92D\u795e\u7ecf\u5143\u57f9\u517b\u7269\u5b66\u4f1a\u4e86\u6253Pong\uff01\u4f7f\u7528\u5b9e\u9a8c\u5ba4\u7684\u4e07\u5143ADC\u8bfb\u53d6\u5c16\u5cf0\u4fe1\u53f7\uff0c\u8bbe\u8ba1\u4e86\u7535\u6781\u5206\u7ec4\u65b9\u6848\uff08\u4f4d\u7f6e\u7f16\u7801+\u8fd0\u52a8\u89e3\u7801+\u5956\u60e9\u53cd\u9988\uff09\uff0c\u7528Cursor\u5199\u4e86Pong\u4ee3\u7801\uff0c\u63a5\u4e0a\u57f9\u517b\u7269\u540e\u76f4\u63a5\u6210\u529f\u4e86\u3002\u5f55\u4e86\u89c6\u9891\u7ed9\u6559\u6388\u770b\uff0c\u4ed6\u5f88\u60ca\u8bb6\uff08\u867d\u7136\u4ed6\u4ee5\u524d\u4e5f\u505a\u8fc7\u7c7b\u4f3c\u7684\u2026\u2026\uff09\n\n"
            "\u5f00\u59cb\u4e3a3D\u7cfb\u7edf\u505a\u51c6\u5907\uff1a\n"
            "- \u8ba2\u8d2d\u4e868\u79cd\u9ad8\u7ea7\u795e\u7ecf\u5143\uff08\u591a\u5df4\u80fa\u795e\u7ecf\u5143\u3001\u76ae\u5c42\u795e\u7ecf\u5143\u3001\u611f\u89c9\u795e\u7ecf\u5143\u7b49\uff09\uff0c\u82b1\u4e86\u6559\u630210\u4e07\u7f8e\u5143\u7684\u5171\u540c\u91c7\u8d2d\n"
            "- \u8bbe\u8ba1\u65b0\u724815\u00d72\u00d72cm\u57f9\u517b\u8154\uff0c\u53d7\u5c0f\u8d3e\u704c\u6d41\u7cfb\u7edf\u542f\u53d1\u6539\u8fdb\u4e86\u6cf5\u7684\u5e03\u5c40\n"
            "- \u542f\u52a8\"keeping_it_alive\"\u9879\u76ee\uff1a\u4e0e\u505a\u4f53\u5916\u80da\u80ce\u7684\u540c\u4e8b\u5c0f\u8d3e\u5408\u4f5c\uff0c\u642d\u5efa\u704c\u6d41\u7cfb\u7edf\u7ef4\u6301\u8111\u7ec4\u7ec7\u957f\u671f\u5b58\u6d3b\n"
            "- \u5f00\u59cb\u9605\u8bfb\u5173\u4e8e\u03bc\u4ecb\u5b50\uff08muon\uff09\u4f5c\u4e3a\u8111\u6210\u50cf\u63a2\u9488\u7684\u7406\u8bba\u6587\u732e"
        )
        a['content_en'] = (
            "Joined the lab, met colleagues (Xiao Jia, Xiao Yi, Xiao Ban), joined Cycling Club and Anime Club.\n"
            "Designed a two-stage experimental system:\n"
            "Stage 1\u2014\"Electrode Testing Kit\": Faraday cage + brain cell culture + side wall electrodes + signal readout, verified by training cultured brain to play Doom.\n"
            "Stage 2\u2014Quantum Signal Verifier: neighboring accelerator lab fires muon beam, Cherenkov detectors read 1-D signal, test whether brain tissue can still function under quantum signals.\n\n"
            "Learned SolidWorks modeling (discovered need for bio-compatible materials), completed bio lab safety training. Grew first brain cells in petri dishes using established recipe, cleanroom engineers fabricated electrode walls in 3 days.\n"
            "Initial discussions with collaborator Dr. Wu about muon beam access: they can reliably shoot muon beams, but since muons are weakly interactive, may need to bring experiment to the accelerator. Estimated 2 months before data collection.\n\n"
            "Major breakthrough: Successfully made 2D neuron culture learn to play Pong! Used the lab's $10k ADC to read spike signals, designed electrode grouping scheme (position encoding + movement decoding + reward/punishment feedback), coded Pong with Cursor, connected to culture and it just worked. Recorded video for professor\u2014he was impressed (though he'd done something similar before...).\n\n"
            "Started preparing for the 3D system:\n"
            "- Ordered 8 types of advanced neurons (dopamine neurons, cortical neurons, sensory neurons, etc.), $100k pooled purchase with neighboring labs\n"
            "- Designed new 15\u00d72\u00d72cm culture chamber, improved pump layout inspired by Jia's perfusion system\n"
            "- Launched \"keeping_it_alive\" project: collaborated with colleague Jia (in vitro embryo specialist) to build perfusion system for long-term brain tissue survival\n"
            "- Started reading theoretical literature on muons as brain imaging probes"
        )
    elif a['month'] == 2:
        a['title'] = "\u7b2c2\u6708\uff1a3D\u7cfb\u7edf\u4e0e\u6cd5\u62c9\u7b2c\u7b3c"
        a['title_en'] = "Month 2: 3D System & Faraday Cage"
        a['content'] = (
            "\u9ad8\u7ea7\u7ec6\u80de\u7c7b\u578b\u5230\u8d27\u3002\u572815\u00d72\u00d72cm\u57f9\u517b\u8154\u4e2d\u7ec4\u88c53D\u8111\u7ec4\u7ec7\uff0c\u914d\u5907\u5c0f\u8d3e\u4f18\u5316\u8fc7\u7684\u704c\u6d41\u7cfb\u7edf\uff08\u6c14\u6ce1\u6355\u96c6\u5668\u3001\u53cc\u5411\u6cf5\u3001\u5206\u79bb\u819c\uff09\u3002\n"
            "\u5b9e\u8df5\u4e2d\u53d1\u73b0\u7ef4\u63013D\u8111\u7ec4\u7ec7\u6bd42D\u590d\u6742\u5f97\u591a\u2014\u2014\u9700\u8981\u7cbe\u786e\u63a7\u5236\u8425\u517b\u6db2\u6d41\u901f\u548c\u6e29\u5ea6\uff0c3D\u7ed3\u6784\u4e2d\u7684\u4fe1\u53f7\u4f20\u5bfc\u4e5f\u66f4\u96be\u9884\u6d4b\u3002\n\n"
            "\u5f00\u59cb\u8bbe\u8ba1\u548c\u642d\u5efa\u4e09\u5c42\u03bc\u91d1\u5c5e\u6cd5\u62c9\u7b2c\u7b3c\uff082m\u00d72m\u00d72m\uff09\uff0c\u76ee\u6807\u5c06\u5185\u90e8\u78c1\u573a\u964d\u4f4e\u5230<1 nT\u3002\u8d2d\u4e70\u6750\u6599\uff0c\u5b66\u4e60\u710a\u63a5\u03bc\u91d1\u5c5e\u677f\u3002\n"
            "\u5c0f\u8d3e\u8bbe\u8ba1\u4e86\u6c14\u6d6e\u9632\u9707\u5e73\u53f0\u6765\u9694\u79bb\u5efa\u7b51\u632f\u52a8\uff08<0.1 \u03bcm\u632f\u5e45\uff09\u3002\u653e\u7f6e\u94a8\u5757\u4f5c\u4e3a\u5df2\u77e5\u6563\u5c04\u53c2\u8003\u7269\uff0c\u7528\u4e8e\u6821\u51c6\u63a2\u6d4b\u5668\u7684\u89d2\u5206\u8fa8\u7387\u548c\u80fd\u91cf\u54cd\u5e94\u3002\n\n"
            "\u7ee7\u7eed\u4e0e\u5434\u535a\u58eb\u56e2\u961f\u8ba8\u8bba\u52a0\u901f\u5668\u5408\u4f5c\u3002\u4ed6\u4eec\u540c\u610f\u63d0\u4f9b\u6bcf\u54688\u5c0f\u65f6\u7684\u03bc\u4ecb\u5b50\u675f\u65f6\uff08200 MeV/c\uff0c\u901a\u91cf~10^6/\u79d2/cm\u00b2\uff09\uff0c\u4f46\u8bbe\u5907\u4f7f\u7528\u8d39\u9ad8\u6602\u3002\u6559\u6388\u5f00\u59cb\u7533\u8bf7\u989d\u5916\u7ecf\u8d39\u3002\n"
            "3D\u57f9\u517b\u7269\u521d\u89c1\u6210\u6548\uff0c\u90e8\u5206\u795e\u7ecf\u5143\u5728\u7535\u6781\u58c1\u4e0a\u6210\u529f\u9644\u7740\u5e76\u663e\u793a\u6d3b\u52a8\u4fe1\u53f7\u3002\n\n"
            "\u5f00\u59cb\u7814\u7a76DIY\u5bb6\u7528\u7092\u83dc\u673a\uff0c\u8ba2\u4e86100\u7f8e\u5143\u7684\u96f6\u4ef6\uff0c\u5f00\u59cb\u8bbe\u8ba1\u673a\u68b0\u7ed3\u6784\u3002"
        )
        a['content_en'] = (
            "Advanced cell types arrived. Assembled 3D brain tissue in the 15\u00d72\u00d72cm culture chamber with Jia's optimized perfusion system (bubble traps, bidirectional pumps, separation membranes).\n"
            "In practice, maintaining 3D brain tissue is much more complex than 2D\u2014requires precise control of nutrient flow rate and temperature, and signal propagation in 3D structures is harder to predict.\n\n"
            "Started designing and building a three-layer mu-metal Faraday cage (2m\u00d72m\u00d72m), targeting internal magnetic field reduction to <1 nT. Purchasing materials, learning to weld mu-metal plates.\n"
            "Jia designed an air-bearing anti-vibration platform to isolate building vibrations (<0.1 \u03bcm amplitude). Placed tungsten blocks as known scattering references for calibrating detector angular resolution and energy response.\n\n"
            "Continued accelerator collaboration discussions with Dr. Wu's team. They agreed to provide 8 hours of muon beam time per week (200 MeV/c, flux ~10^6/sec/cm\u00b2), but equipment usage fees were high. Professor started applying for additional funding.\n"
            "3D culture showing initial results\u2014some neurons successfully attached to electrode walls and displaying activity signals.\n\n"
            "Started researching DIY home cooking machine, ordered $100 worth of parts, began designing the mechanical structure."
        )

# ============================================================
# UPDATE SOCIAL ANCHORS (Month 1) for consistency
# ============================================================

social = data['kurisu']['story']['socialAnchors']
for s in social:
    if s['month'] == 1:
        s['content'] = "\u521a\u5165\u804c\uff0c\u8dd1\u904d\u5404\u4e2a\u5b9e\u9a8c\u5ba4\u4e86\u89e3\u5927\u5bb6\u7684\u5206\u5de5\u548c\u6027\u683c\u3002\u5f88\u5feb\u8ddf\u5c0f\u8d3e\u3001\u5c0f\u4f0a\u3001\u5c0f\u6591\u6df7\u719f\u4e86\u3002\u52a0\u5165\u4e86\u9a91\u884c\u793e\u548c\u52a8\u6f2b\u793e\uff0c\u8ba4\u8bc6\u4e86Emily\u3002\u6bcf\u5468\u4e94\u8ddf\u7740\u4e01\u6559\u6388\u53bb\u9152\u5427\u559d\u9152\u804a\u5929\u3002\u5f00\u59cb\u770b\u300a\u659a\uff01\u8d64\u7ea2\u4e4b\u779b\u300b\uff0c\u88ab\u5e0c\u5c14\u7684\u6545\u4e8b\u6df1\u6df1\u89e6\u52a8\u3002\u611f\u89c9\u8ddf\u5b9e\u9a8c\u5ba4\u7684\u4eba\u867d\u7136\u5f88\u597d\uff0c\u4f46\u6709\u79cd\u8bf4\u4e0d\u51fa\u7684\u8ddd\u79bb\u611f\u3002\u5f00\u59cb\u60f3\u5728\u7f51\u4e0a\u627e\u4eba\u804a\u5929\u3002"
        s['content_en'] = "Just joined, visited all the labs to understand everyone's roles and personalities. Quickly became close with Xiao Jia, Xiao Yi, and Xiao Ban. Joined the Cycling Club and Anime Club, met Emily. Started going to the bar with Professor Ding every Friday. Started watching Akame ga Kill, deeply moved by Sheele's story. Felt like while the people at the lab are great, there's an inexplicable sense of distance. Started wanting to find people to talk to online."

# ============================================================
# Write the updated data back (safe write: temp file first)
# ============================================================
prod_path = r'C:\Users\user\Desktop\data_labeler\data\character_profiles.json'
testing_path = r'C:\Users\user\Desktop\data_labeler_testing_env\data\character_profiles.json'

# First serialize to string to validate before touching any file
output = json.dumps(data, ensure_ascii=False, indent=2)
print(f"Serialized OK: {len(output)} chars")

# Write production
with open(prod_path, 'w', encoding='utf-8') as f:
    f.write(output)

print("Production character_profiles.json updated successfully!")
print(f"Kurisu now has {len(data['kurisu']['schedules'])} schedule entries:")
for s in data['kurisu']['schedules']:
    print(f"  Day {s['day']}: {s['title_en']}")

# Also update the testing env copy
with open(testing_path, 'w', encoding='utf-8') as f:
    f.write(output)
print(f"\nTesting env copy also updated.")

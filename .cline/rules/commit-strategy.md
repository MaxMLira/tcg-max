# Commit strategy

All the work we do here must be split into multiple commits. The main idea is for these commits to be well delimited and isolated from each other whenever possible. I want to make code review easier and to make commit reversions a good tool for corrections.
Sometimes, instead of writing a correction that runs across code changed in many recent commits, the idea would be to revert a single commit and do the correction over code that that commit would be referencing.
I know it is sometimes difficult to foresee such things, but with good practices we can get most of them correct.

## Rule 1: don't mix frontend and backend code in a commit
Sometimes our solutions require changes in the backend, in the frontend and even in the proto files. What should happen:
- Backend changes must have their own commits
- Frontend changes must have their own commits
- Proto file changes must have their own commits
- Documentation updates must have their own commits
The commits cannot have a mix of these types of change. They will be evaluated and reviewed by different people.

## Rule 2: commits cannot be too long or too short
Long commits are bad because they hardly can be reverted. But if they are too short, we will have an inconvenient amount of commits to deal with. The main idea is that, when we are implementing a plan, each topic of the plan must have at least one commit.
When the topic requires many actions (e.g.: create functions to do something, correct queries to include a new parameter, change model structs to embed the new field, update auth logic...), then each action (here each action is being represented by a sentence that starts with a verb) must have its own commit.

## Rule 3: commits with title and justifications
Every commit must have a message with a short title indicating what was done and a text that explain what was done in that commit and why it was done. This text must include only important information and must be short enough (maximum 10 lines) to ease the review process.
Most importantly: what was done can be described very shortly, because the user can see the code changes. But the reason why it was done must be **very** clear and reasonable.

## Rule 4: commits must have a sequence
The commits must be done in a logic sequence. We want to avoid situations such as implementing something that has a dependency in one commit and implementing the dependency in a future commit.
Whenever possible (though we know it might not always be possible), every commit must be created so that the code up to that commit has no syntax errors or missing dependencies and builds properly. This also affects the commit sequence.

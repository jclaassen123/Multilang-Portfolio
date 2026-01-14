import random

class Wordle:
    def __init__(self, word_list=None):
        """
        Initialize the Wordle game logic.

        Parameters:
        - word_list: Optional list of 5-letter words to choose from.
                     If None, uses the default list.
        """
        # Use provided word list or default list
        self.word_list = word_list or [
            "apple", "grape", "mango", "peach", "berry",
            "lemon", "melon", "chess", "flame", "brick",
            "train", "plant", "shore", "stone", "snake",
            "lunar", "tiger", "plane", "quake", "glove",
            "drain", "plaza", "crown", "slice", "pride",
            "ghost", "blaze", "frost", "shark", "vigor",
            "alien", "angel", "beach", "baker", "candy",
            "cabin", "dance", "delta", "eagle", "fable",
            "faith", "giant", "globe", "harpy", "honey",
            "ivory", "joker", "knife", "leech", "light",
            "magic", "mimic", "noble", "nymph", "ocean",
            "olive", "pixel", "queen", "quilt", "radio",
            "robot", "salsa", "scout", "tango", "tempo",
            "umbra", "uncle", "vocal", "vivid", "watch",
            "waltz", "quick", "brown", "yacht", "yield",
            "zebra", "zesty", "acorn", "brave", "caper",
            "daisy", "ember", "frost", "glint", "haunt",
            "inbox", "jolly", "karma", "latch", "mirth",
            "nudge", "orbit", "plume", "quark", "raven",
            "swirl", "trove", "ultra", "vixen", "wreak",
            "yodel", "night", "amber", "bliss", "crisp"
        ]
        # Initialize the game state
        self.reset_game()

    def reset_game(self):
        """
        Start a new game:
        - Randomly select a target word
        - Clear previous attempts
        - Set maximum allowed attempts
        """
        self.target_word = random.choice(self.word_list)  # Randomly pick target word
        self.attempts = []                                # Clear attempts history
        self.max_attempts = 6                             # Standard Wordle has 6 attempts
        print(f"Target word: {self.target_word}") # Optional debug output

    def guess(self, word):
        """
        Process a guess from the player.

        Parameters:
        - word: 5-letter string guessed by the player

        Returns:
        - A string of 5 symbols representing feedback:
          🟩 = correct letter, correct position
          🟨 = correct letter, wrong position
          ⬜ = letter not in word
        """
        word = word.lower()
        if len(word) != 5:
            return "Word must be 5 letters!"  # Input validation

        feedback = ""
        for i in range(5):
            if word[i] == self.target_word[i]:
                feedback += "🟩"  # Correct letter, correct spot
            elif word[i] in self.target_word:
                feedback += "🟨"  # Correct letter, wrong spot
            else:
                feedback += "⬜"  # Letter not in the target word

        # Record this guess and its feedback
        self.attempts.append((word, feedback))
        return feedback

    def is_game_over(self):
        """
        Check if the game has ended, either by win or reaching max attempts.

        Returns:
        - Tuple (over: bool, message: str)
          - over: True if the game has ended
          - message: "You won!" if player guessed correctly,
                     "You lost! The word was '...'" if max attempts reached,
                     "" otherwise.
        """
        # Player wins if last guess matches target word
        if self.attempts and self.attempts[-1][0] == self.target_word:
            return True, "You won!"

        # Player loses if maximum attempts reached
        if len(self.attempts) >= self.max_attempts:
            return True, f"You lost! The word was '{self.target_word}'"

        # Game continues
        return False, ""

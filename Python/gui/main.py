import tkinter as tk
from tkinter import messagebox
from wordle import Wordle

class WordleApp:
    """
    GUI application for the Wordle game using Tkinter.

    Features:
    - 6x5 grid for guesses
    - On-screen clickable keyboard
    - Backspace and Submit buttons
    - Keyboard colors updated based on feedback
    """

    COLORS = {
        "default": "#d3d6da",  # Light gray for unused grid boxes
        "correct": "#6aaa64",  # Green for correct letters in correct position
        "present": "#c9b458",  # Yellow for correct letters in wrong position
        "absent": "#595959"    # Dark gray for keyboard letters not in word
    }

    def __init__(self, root: tk.Tk):
        """
        Initialize the Wordle GUI application.

        Args:
            root (tk.Tk): Root Tkinter window.
        """
        self.root = root
        self.root.title("Wordle GUI")
        self.game = Wordle()  # Wordle game logic object
        self.letter_labels = {}  # Stores keyboard letter labels
        self.grid_labels = []    # 6x5 grid for guesses
        self.current_row = 0     # Track which row of the grid we're filling
        self.current_col = 0     # Track current letter position in row

        # Create GUI elements
        self.create_widgets()

    def create_widgets(self):
        """
        Create all GUI components:
        - Grid for guesses
        - On-screen keyboard
        - Backspace, Submit, Reset, Quit buttons
        """
        # --- Word grid ---
        self.grid_frame = tk.Frame(self.root)
        self.grid_frame.pack(pady=10)

        # Create 6x5 label grid
        for r in range(6):
            row_labels = []
            for c in range(5):
                lbl = tk.Label(
                    self.grid_frame,
                    text="",
                    width=4,
                    height=2,
                    bg=self.COLORS["default"],
                    relief="solid",
                    font=("Consolas", 24),
                    borderwidth=2,
                    fg="black"
                )
                lbl.grid(row=r, column=c, padx=5, pady=5)
                row_labels.append(lbl)
            self.grid_labels.append(row_labels)

        # --- On-screen keyboard ---
        keyboard_frame = tk.Frame(self.root)
        keyboard_frame.pack(pady=10)

        rows = ["QWERTYUIOP", "ASDFGHJKL", "ZXCVBNM"]
        for r, row in enumerate(rows):
            row_frame = tk.Frame(keyboard_frame)
            row_frame.pack(pady=2)
            for letter in row:
                lbl = tk.Label(
                    row_frame,
                    text=letter,
                    width=4,
                    height=2,
                    bg=self.COLORS["default"],
                    relief="raised",
                    fg="black",
                    font=("Consolas", 12)
                )
                lbl.pack(side="left", padx=2, pady=2)
                # Bind click to add_letter function
                lbl.bind("<Button-1>", lambda e, l=letter: self.add_letter(l))
                self.letter_labels[letter] = lbl

        # --- Backspace button ---
        backspace = tk.Label(
            row_frame,
            text="⌫",
            width=4,
            height=2,
            bg="#bbbbbb",
            relief="raised",
            fg="black",
            font=("Consolas", 12)
        )
        backspace.pack(side="left", padx=2, pady=2)
        backspace.bind("<Button-1>", lambda e: self.backspace_letter())
        self.letter_labels["BACKSPACE"] = backspace

        # --- Submit button ---
        self.submit_button = tk.Button(self.root, text="Submit", command=self.make_guess)
        self.submit_button.pack(pady=5)

        # --- Reset / Quit buttons ---
        self.reset_button = tk.Button(self.root, text="New Game", command=self.reset_game)
        self.reset_button.pack(pady=5)
        self.quit_button = tk.Button(self.root, text="Quit", command=self.root.destroy)
        self.quit_button.pack(pady=5)

    def add_letter(self, letter: str):
        """
        Add a letter from the on-screen keyboard to the current row.
        Only adds if there are less than 5 letters in the row.

        Args:
            letter (str): Letter to add.
        """
        if self.current_col < 5 and self.current_row < 6:
            lbl = self.grid_labels[self.current_row][self.current_col]
            lbl.config(text=letter)
            self.current_col += 1

    def backspace_letter(self):
        """
        Remove the last letter from the current row.
        """
        if self.current_col > 0:
            self.current_col -= 1
            lbl = self.grid_labels[self.current_row][self.current_col]
            lbl.config(text="")

    def make_guess(self):
        """
        Submit the current row as a guess:
        - Validate the guess length
        - Update the grid colors and keyboard colors
        - Check for win/loss after updating the GUI
        """
        if self.current_col != 5:
            messagebox.showinfo("Invalid Guess", "Please enter 5 letters!")
            return

        # Build the word from the current row
        word = "".join(lbl.cget("text") for lbl in self.grid_labels[self.current_row]).upper()
        feedback = self.game.guess(word.lower())

        # --- Update the grid for this row ---
        for i in range(5):
            lbl = self.grid_labels[self.current_row][i]
            if feedback[i] == "🟩":
                lbl.config(bg=self.COLORS["correct"])
            elif feedback[i] == "🟨":
                lbl.config(bg=self.COLORS["present"])
            else:
                lbl.config(bg=self.COLORS["default"])

        # --- Update keyboard colors ---
        self.update_keyboard(word, feedback)

        # Force GUI to redraw before showing popup
        self.root.update_idletasks()

        # Move to next row
        self.current_row += 1
        self.current_col = 0

        # --- Check if game over after updating GUI ---
        over, msg = self.game.is_game_over()
        if over:
            messagebox.showinfo("Game Over", msg)

    def update_keyboard(self, word: str, feedback: str):
        """
        Update the colors of the on-screen keyboard based on feedback.

        Args:
            word (str): The guessed word.
            feedback (str): Feedback string (🟩, 🟨, ⬜).
        """
        for i, letter in enumerate(word):
            lbl = self.letter_labels.get(letter.upper())
            if not lbl:
                continue
            if feedback[i] == "🟩":
                lbl.config(bg=self.COLORS["correct"])
            elif feedback[i] == "🟨":
                if lbl.cget("bg") != self.COLORS["correct"]:
                    lbl.config(bg=self.COLORS["present"])
            else:
                if lbl.cget("bg") not in [self.COLORS["correct"], self.COLORS["present"]]:
                    lbl.config(bg=self.COLORS["absent"])

    def reset_game(self):
        """
        Reset the game to initial state:
        - Clear the grid
        - Reset keyboard colors
        - Reset game logic
        """
        self.game.reset_game()
        self.current_row = 0
        self.current_col = 0

        # Reset grid
        for r in range(6):
            for c in range(5):
                lbl = self.grid_labels[r][c]
                lbl.config(text="", bg=self.COLORS["default"], fg="black")

        # Reset keyboard
        for lbl in self.letter_labels.values():
            if lbl.cget("text") == "⌫":
                lbl.config(bg="#bbbbbb")
            else:
                lbl.config(bg=self.COLORS["default"])


if __name__ == "__main__":
    # Launch the Wordle GUI
    root = tk.Tk()
    app = WordleApp(root)
    root.mainloop()

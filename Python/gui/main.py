import tkinter as tk
from tkinter import messagebox
from wordle import Wordle

class WordleApp:
    # Define colors for grid and keyboard
    COLORS = {
        "default": "#d3d6da",       # Light gray for unused grid boxes
        "correct": "#6aaa64",       # Green for correct letters in correct position
        "present": "#c9b458",       # Yellow for correct letters in wrong position
        "absent": "#595959"         # Dark gray for keyboard letters not in word
    }

    def __init__(self, root):
        """
        Initialize the Wordle GUI application.
        """
        self.root = root
        self.root.title("Wordle GUI")       # Window title
        self.game = Wordle()                # Wordle game logic object

        self.letter_labels = {}             # Stores keyboard letter labels
        self.grid_labels = []               # 6x5 grid for guesses
        self.current_row = 0                # Track which row of the grid we're filling

        self.create_widgets()               # Setup all GUI components

    def create_widgets(self):
        """
        Create all GUI components: grid, typing box, keyboard, reset and quit buttons.
        """
        # --- Word grid at the top ---
        self.grid_frame = tk.Frame(self.root)
        self.grid_frame.pack(pady=10)

        for r in range(6):  # 6 rows for 6 guesses
            row_labels = []
            for c in range(5):  # 5 columns for 5 letters
                lbl = tk.Label(
                    self.grid_frame,
                    text="",
                    width=4,
                    height=2,
                    bg=self.COLORS["default"],
                    relief="solid",
                    font=("Consolas", 24),
                    borderwidth=2,
                    fg="black"  # Letters are black
                )
                lbl.grid(row=r, column=c, padx=5, pady=5)  # Place in grid
                row_labels.append(lbl)
            self.grid_labels.append(row_labels)

        # --- Typing box above Guess button ---
        input_frame = tk.Frame(self.root)
        input_frame.pack(pady=5)

        self.guess_entry = tk.Entry(input_frame, font=("Consolas", 16), justify="center")
        self.guess_entry.pack(side="top", pady=2)
        self.guess_entry.bind("<Return>", lambda e: self.make_guess())  # Press Enter to guess

        self.guess_button = tk.Button(input_frame, text="Guess", command=self.make_guess)
        self.guess_button.pack(side="top", pady=2)

        # --- On-screen keyboard below ---
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
                    bg="#d3d6da",
                    relief="raised",
                    fg="black",
                    font=("Consolas", 12)
                )
                lbl.pack(side="left", padx=2, pady=2)
                lbl.bind("<Button-1>", lambda e, l=letter: self.add_letter(l))  # Click letter to type
                self.letter_labels[letter] = lbl

        # --- Backspace button at the end of last row ---
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
        backspace.bind("<Button-1>", lambda e: self.backspace_letter())  # Delete last letter
        self.letter_labels["BACKSPACE"] = backspace

        # --- Reset button below keyboard ---
        self.reset_button = tk.Button(self.root, text="New Game", command=self.reset_game)
        self.reset_button.pack(pady=5)

        # --- Quit button below Reset button ---
        self.quit_button = tk.Button(self.root, text="Quit", command=self.root.destroy)
        self.quit_button.pack(pady=5)

    def add_letter(self, letter):
        """
        Add a letter from the on-screen keyboard to the typing box.
        Only allows up to 5 letters.
        """
        if len(self.guess_entry.get()) < 5:
            self.guess_entry.insert(tk.END, letter)

    def backspace_letter(self):
        """
        Remove the last letter from the typing box.
        """
        current = self.guess_entry.get()
        if current:
            self.guess_entry.delete(len(current)-1, tk.END)

    def make_guess(self):
        """
        Handles when the user submits a guess.
        Updates the grid, keyboard, and checks for win/loss.
        """
        word = self.guess_entry.get().strip().upper()
        if len(word) != 5:
            messagebox.showinfo("Invalid Guess", "Word must be exactly 5 letters!")
            return

        feedback = self.game.guess(word.lower())  # Get feedback from Wordle logic
        self.guess_entry.delete(0, tk.END)

        # --- Update the grid for current row ---
        for i in range(5):
            lbl = self.grid_labels[self.current_row][i]
            lbl.config(text=word[i], fg="black")  # Always black letters
            if feedback[i] == "🟩":
                lbl.config(bg=self.COLORS["correct"])
            elif feedback[i] == "🟨":
                lbl.config(bg=self.COLORS["present"])
            else:
                lbl.config(bg=self.COLORS["default"])  # Keep light gray for absent letters

        # --- Update keyboard colors ---
        self.update_keyboard(word, feedback)

        # Force GUI to redraw so the final row shows before popup
        self.root.update_idletasks()

        # Move to the next row
        self.current_row += 1

        # --- Check if game over (win or loss) ---
        over, msg = self.game.is_game_over()
        if over:
            messagebox.showinfo("Game Over", msg)

    def update_keyboard(self, word, feedback):
        """
        Update the colors of the on-screen keyboard based on feedback.
        Green = correct, Yellow = present, Dark Gray = absent.
        """
        for i, letter in enumerate(word):
            lbl = self.letter_labels.get(letter.upper())
            if not lbl:
                continue
            if feedback[i] == "🟩":
                lbl.config(bg=self.COLORS["correct"])
            elif feedback[i] == "🟨":
                if lbl.cget("bg") != self.COLORS["correct"]:  # Don't overwrite green
                    lbl.config(bg=self.COLORS["present"])
            else:  # Letter not in word → dark gray
                if lbl.cget("bg") not in [self.COLORS["correct"], self.COLORS["present"]]:
                    lbl.config(bg=self.COLORS["absent"])

    def reset_game(self):
        """
        Resets the Wordle game, clears the grid and keyboard colors.
        """
        self.game.reset_game()
        self.current_row = 0

        # Reset grid labels
        for r in range(6):
            for c in range(5):
                lbl = self.grid_labels[r][c]
                lbl.config(text="", bg=self.COLORS["default"], fg="black")

        # Reset keyboard colors
        for lbl in self.letter_labels.values():
            if lbl.cget("text") == "⌫":
                lbl.config(bg="#bbbbbb")
            else:
                lbl.config(bg="#d3d6da")

if __name__ == "__main__":
    root = tk.Tk()
    app = WordleApp(root)
    root.mainloop()

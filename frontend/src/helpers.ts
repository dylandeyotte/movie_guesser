export type status = "hidden" | "correct" | "failed";

export type gameInfo = {
  actor: string;
  gamedate: string;
};

export type filmCard = {
  title: string;
  poster: string;
  status: status;
};

export type FilmAnswers = {
  film1: string;
  film1_poster: string;
  film2: string;
  film2_poster: string;
  film3: string;
  film3_poster: string;
};

export type guessResponse = {
  verdict: boolean;
  film_number: number;
  strikes: number;
  guess: string;
  playerid: string;
  repeat: boolean;
  poster_path: string[];
};

export type guessInfo = {
  Date: string;
  FilmNumber: number;
  Guess: string;
  Verdict: boolean;
};

export type gameState = {
  date: string;
  actor: string;
  playerid: string;
  guesses: guessInfo[];
  strikes: number;
  posters: string[];
  game_over: boolean;
  answers: FilmAnswers;
};

import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { CircleQuestionMark, ChartNoAxesColumn, X } from "lucide-react";

import type { status, gameInfo, filmCard, guessResponse, gameState } from "./helpers";

type gameEnd = {
  victory: boolean;
  defeat: boolean;
};

type searchResults = {
  title: string;
  id: number;
};

export function Home() {
  const navigate = useNavigate();
  const [info, setInfo] = useState<gameInfo>();
  const [guess, setGuess] = useState("");
  const [results, setResults] = useState<searchResults[]>([]);
  const [HTPDisplay, setHTPDisplay] = useState(false);
  const [victoryText, setVictoryText] = useState("");
  const [gameEnd, setGameEnd] = useState<gameEnd>();
  const [filmOne, setFilmOne] = useState<filmCard>();
  const [filmTwo, setFilmTwo] = useState<filmCard>();
  const [filmThree, setFilmThree] = useState<filmCard>();
  const [gameState, setGameState] = useState<gameState>();
  const [guessResponse, setguessResponse] = useState<guessResponse>();
  const [incorrectGuess, setIncorrectGuess] = useState<Set<string>>(new Set());
  const guessInputRef = useRef<HTMLInputElement>(null);

  const playerID = localStorage.getItem("playerID") ?? crypto.randomUUID();

  localStorage.setItem("playerID", playerID);

  function setFilmHelper(film: string, path: string, status: status, filmNumber: number) {
    switch (filmNumber) {
      case 1:
        setFilmOne({
          title: film,
          poster: `https://image.tmdb.org/t/p/w500${path}`,
          status: status,
        });
        break;
      case 2:
        setFilmTwo({
          title: film,
          poster: `https://image.tmdb.org/t/p/w500${path}`,
          status: status,
        });
        break;
      case 3:
        setFilmThree({
          title: film,
          poster: `https://image.tmdb.org/t/p/w500${path}`,
          status: status,
        });
        break;
    }
  }

  async function gameStatePull(date: string) {
    try {
      const response = await fetch("http://localhost:8080/api/gamestate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-Player-ID": playerID,
        },
        body: JSON.stringify({
          date: date,
        }),
      });
      const data = await response.json();

      // Set victory text
      data?.strikes && data.strikes > 0 ? setVictoryText("You did it!") : setVictoryText("Perfect!");

      // End game if victorious
      if (data.victory === true) {
        setGameEnd({
          victory: true,
          defeat: false,
        });
      }
      // End game if failed
      if (data.game_over === true) {
        setGameEnd({
          victory: false,
          defeat: true,
        });
        // Reveal answers
        setFilmHelper(data.answers.film1, data.answers.film1_poster, "failed", 1);
        setFilmHelper(data.answers.film2, data.answers.film2_poster, "failed", 2);
        setFilmHelper(data.answers.film3, data.answers.film3_poster, "failed", 3);
      }
      // Display guessed films
      for (const guess of data.guesses) {
        if (guess.Verdict === true) {
          switch (guess.FilmNumber) {
            case 1:
              setFilmHelper(guess.Guess, data.posters[0], "correct", 1);
              break;
            case 2:
              setFilmHelper(guess.Guess, data.posters[1], "correct", 2);
              break;
            case 3:
              setFilmHelper(guess.Guess, data.posters[2], "correct", 3);
              break;
          }
        }
      }

      setGameState(data);

      // Display incorrect guesses
      if (data) {
        for (const guess of data?.guesses) {
          guess.Verdict === false && setIncorrectGuess((prev) => new Set(prev).add(guess.Guess));
        }
      }
    } catch (err) {
      console.error(err);
      throw err;
    }
  }

  async function gamePull() {
    try {
      const response = await fetch("http://localhost:8080/api/actor", {
        headers: {
          "Content-Type": "application/json",
        },
      });
      const data = await response.json();
      setInfo(data);
      return data;
    } catch (err) {
      console.error(err);
      throw err;
    }
  }

  const submitGuess = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();

    guessHelper(guess, false);
  };

  const submitGiveUp = async () => {
    guessHelper("", true);
  };

  const guessHelper = async (guess: string, giveUp: boolean) => {
    const response = await fetch("http://localhost:8080/api/guess", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Player-ID": playerID,
      },
      body: JSON.stringify({
        guess: guess,
        gamedate: info?.gamedate,
        giveup: giveUp,
      }),
    });
    const data = await response.json();

    // Set victory text
    data?.strikes && data.strikes > 0 ? setVictoryText("You did it!") : setVictoryText("Perfect!");

    // End game if victorious
    if (info?.gamedate) {
      if (data.game_over === true) {
        setGameEnd({
          victory: false,
          defeat: true,
        });
        await gameStatePull(info?.gamedate);
      }
      // End game if failed
      if (data.victory === true) {
        setGameEnd({
          victory: true,
          defeat: false,
        });
        await gameStatePull(info?.gamedate);
      }
    }

    // Display correct guess
    switch (data.film_number) {
      case 1:
        setFilmHelper(data.guess, data.poster_path[0], data.verdict === true ? "correct" : "hidden", 1);
        break;
      case 2:
        setFilmHelper(data.guess, data.poster_path[0], data.verdict === true ? "correct" : "hidden", 2);
        break;
      case 3:
        setFilmHelper(data.guess, data.poster_path[0], data.verdict === true ? "correct" : "hidden", 3);
        break;
    }
    // Display incorrect guess
    if (data.verdict == false) {
      setIncorrectGuess(new Set(incorrectGuess).add(guess));
    }
    setguessResponse(data);

    // Reset guess bar
    if (response.ok) {
      setGuess("");
    }
  };

  async function searchFilms(title: string) {
    {
      try {
        const response = await fetch("http://localhost:8080/api/search", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            film: title,
            date: info?.gamedate,
          }),
        });

        const data = await response.json();
        return data;
      } catch (err) {
        console.error(err);
        throw err;
      }
    }
  }

  useEffect(() => {
    async function loadGame() {
      const gameInfo = await gamePull();
      await gameStatePull(gameInfo.gamedate);
    }
    loadGame();
  }, []);

  useEffect(() => {
    if (guess.trim().length < 2) {
      setResults([]);
      return;
    }

    const timer = setTimeout(async () => {
      const data = await searchFilms(guess);
      setResults(data);
    }, 300);

    return () => clearTimeout(timer);
  }, [guess]);

  console.log(results);

  return (
    <div className="background">
      <div className="title-row">
        <div className="actor-name">{info?.actor}</div>
        <div className="home-buttons">
          <button className="HTP-button" onClick={() => setHTPDisplay(true)}>
            <CircleQuestionMark />
          </button>
          <button className="stats-button" onClick={() => navigate("/stats")}>
            <ChartNoAxesColumn />
          </button>
        </div>
        {HTPDisplay && (
          <div className="modal-overlay">
            <div className="modal">
              <button className="x-button" onClick={() => setHTPDisplay(false)}>
                <X size={20} />
              </button>
              <h2>How to Play</h2>
              <p>Guess the three most popular movies for today's actor as determined by TMDB.</p>
              <p>Each incorrect guess gives you a strike.</p>
              <p>Three strikes and the game is over.</p>
            </div>
          </div>
        )}
      </div>
      <div className="poster-container">
        <div className="poster-card">
          <div className={`poster ${filmOne?.status}`}>
            <img src={filmOne?.poster} className="poster-img" />
          </div>
          <div className="film-title">{filmOne?.title ? filmOne.title : "???"}</div>
        </div>
        <div className="poster-card">
          <div className={`poster ${filmTwo?.status}`}>
            <img src={filmTwo?.poster} className="poster-img" />
          </div>
          <div className="film-title">{filmTwo?.title ? filmTwo.title : "???"}</div>
        </div>
        <div className="poster-card">
          <div className={`poster ${filmThree?.status}`}>
            <img src={filmThree?.poster} className="poster-img" />
          </div>
          <div className="film-title">{filmThree?.title ? filmThree.title : "???"}</div>
        </div>
      </div>
      <div>
        {gameEnd?.victory === true ? (
          <div className="end-text">{victoryText}</div>
        ) : gameEnd?.defeat === true ? (
          <div className="end-text">Game Over</div>
        ) : (
          <div className="glorp">
            <div className="guess-search">
              {results.length > 0 && (
                <div className="search-results">
                  {results.map((movie) => (
                    <button
                      className="search-results-button"
                      type="button"
                      key={movie.id}
                      onClick={() => {
                        setGuess(movie.title);
                        guessInputRef.current?.focus();
                      }}
                    >
                      {movie.title}
                    </button>
                  ))}
                </div>
              )}
              <div className="guess-bar-and-button">
                <form className="guess-bar-container" onSubmit={submitGuess}>
                  <input
                    className="guess-bar"
                    type="guess"
                    ref={guessInputRef}
                    value={guess}
                    onChange={(e) => setGuess(e.target.value)}
                    placeholder="Film"
                  />
                </form>
                <div className="give-up-container">
                  <button className="give-up-button" onClick={submitGiveUp}>
                    Give Up
                  </button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
      <div className="strikes">
        {(guessResponse?.strikes ?? 0) >= 1 || (gameState?.strikes ?? 0) >= 1 ? (
          <span className="strike-box"></span>
        ) : (
          <span className="strike-box-empty"></span>
        )}
        {(guessResponse?.strikes ?? 0) >= 2 || (gameState?.strikes ?? 0) >= 2 ? (
          <span className="strike-box"></span>
        ) : (
          <span className="strike-box-empty"></span>
        )}
        {guessResponse?.strikes === 3 || gameState?.strikes === 3 ? (
          <span className="strike-box"></span>
        ) : (
          <span className="strike-box-empty"></span>
        )}
      </div>
      <div className="missed-guess">
        {[...incorrectGuess].map((guess) => (
          <div key={guess}>{guess}</div>
        ))}
      </div>
    </div>
  );
}

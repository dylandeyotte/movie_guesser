import { useState } from "react";
import { useEffect } from "react";

import type { status, gameInfo, filmCard, guessResponse, gameState } from "./helpers";

type gameEnd = {
  victory: boolean;
  defeat: boolean;
};

export function Home() {
  const [info, setInfo] = useState<gameInfo>();
  const [guess, setGuess] = useState("");
  const [gameEnd, setGameEnd] = useState<gameEnd>();
  const [filmOne, setFilmOne] = useState<filmCard>();
  const [filmTwo, setFilmTwo] = useState<filmCard>();
  const [filmThree, setFilmThree] = useState<filmCard>();
  const [gameState, setGameState] = useState<gameState>();
  const [guessResponse, setguessResponse] = useState<guessResponse>();
  const [incorrectGuess, setIncorrectGuess] = useState<Set<string>>(new Set());

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

  async function gameStatePull() {
    try {
      const response = await fetch("http://localhost:8080/api/gamestate", {
        headers: {
          "Content-Type": "application/json",
          "X-Player-ID": playerID,
        },
      });
      const data = await response.json();

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
          console.log(guess.Guess);
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
    } catch (err) {
      console.error(err);
      throw err;
    }
  }

  const submitGuess = async (e: React.SyntheticEvent<HTMLFormElement>) => {
    e.preventDefault();

    const response = await fetch("http://localhost:8080/api/guess", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Player-ID": playerID,
      },
      body: JSON.stringify({
        guess: guess,
        gamedate: info?.gamedate,
      }),
    });

    const data = await response.json();

    // End game if victorious
    if (data.game_over === true) {
      setGameEnd({
        victory: false,
        defeat: true,
      });
      await gameStatePull(); // MIDNIGHT OVERFLOW
    }
    // End game if failed
    if (data.victory === true) {
      setGameEnd({
        victory: true,
        defeat: false,
      });
    }
    console.log(data);

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

  useEffect(() => {
    gamePull();
    gameStatePull();
  }, []);

  return (
    <div className="background">
      <div className="actor-name">{info?.actor}</div>
      <div className="poster-container">
        <div className="poster-card">
          <div className={`poster ${filmOne?.status}`}>
            <img src={filmOne?.poster} className="poster-img" />
          </div>
          <div className="title">{filmOne?.title ? filmOne.title : "???"}</div>
        </div>
        <div className="poster-card">
          <div className={`poster ${filmTwo?.status}`}>
            <img src={filmTwo?.poster} className="poster-img" />
          </div>
          <div className="title">{filmTwo?.title ? filmTwo.title : "???"}</div>
        </div>
        <div className="poster-card">
          <div className={`poster ${filmThree?.status}`}>
            <img src={filmThree?.poster} className="poster-img" />
          </div>
          <div className="title">{filmThree?.title ? filmThree.title : "???"}</div>
        </div>
      </div>
      <div>
        {gameEnd?.victory === true ? (
          <div className="end-text">You did it!</div>
        ) : gameEnd?.defeat === true ? (
          <div className="end-text">Game Over</div>
        ) : (
          <div>
            <form className="guess-box" onSubmit={submitGuess}>
              <input className="box" type="guess" value={guess} onChange={(e) => setGuess(e.target.value)} placeholder="Film" />
            </form>
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
          <div>{guess}</div>
        ))}
      </div>
    </div>
  );
}

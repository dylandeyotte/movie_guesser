import { useState } from "react";
import { useEffect } from "react";

import type { status, gameInfo, filmCard, guessResponse, gameState } from "./helpers";

export function Home() {
  const [info, setInfo] = useState<gameInfo>();
  const [guess, setGuess] = useState("");
  const [gameOver, setGameOver] = useState(Boolean);
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

      if (data.game_over === true) {
        setGameOver(true);
        setFilmHelper(data.answers.film1, data.answers.film1_poster, "failed", 1);
        setFilmHelper(data.answers.film2, data.answers.film2_poster, "failed", 2);
        setFilmHelper(data.answers.film3, data.answers.film3_poster, "failed", 3);
      }

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
      console.log(data.posters);

      setGameState(data);

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

    if (data.game_over === true) {
      await gameStatePull();
    }
    console.log(data);

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

    if (data.verdict == false) {
      setIncorrectGuess(new Set(incorrectGuess).add(guess));
    }
    setguessResponse(data);

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
        {gameOver ? (
          "Game over"
        ) : (
          <div>
            <form className="guess-box" onSubmit={submitGuess}>
              <input className="box" type="guess" value={guess} onChange={(e) => setGuess(e.target.value)} placeholder="Film" />
            </form>
            <div>
              {gameState?.guesses?.map((guess) => (
                <div>{guess.Verdict === false && `${guess.Guess}: incorrect`}</div>
              ))}
            </div>
            <div>
              {[...incorrectGuess].map((guess) => (
                <div>{guess}: incorrect</div>
              ))}
            </div>
            <div>strikes: {guessResponse?.strikes ?? gameState?.strikes}</div>
          </div>
        )}
      </div>
    </div>
  );
}

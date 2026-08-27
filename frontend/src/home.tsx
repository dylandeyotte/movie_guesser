import { useState } from "react";
import { useEffect } from "react";

type gameInfo = {
  actor: string;
  gamedate: string;
};

type guessResponse = {
  verdict: boolean;
  film_number: number;
  strikes: number;
  guess: string;
  playerid: string;
  repeat: boolean;
  poster_path: string[];
};

type guessInfo = {
  Date: string;
  FilmNumber: number;
  Guess: string;
  Verdict: boolean;
};

type gameState = {
  date: string;
  actor: string;
  playerid: string;
  guesses: guessInfo[];
  strikes: number;
  posters: string[];
};

export function Home() {
  const [info, setInfo] = useState<gameInfo>();
  const [guess, setGuess] = useState("");
  const [filmOne, setFilmOne] = useState("");
  const [filmTwo, setFilmTwo] = useState("");
  const [filmThree, setFilmThree] = useState("");
  const [filmOnePoster, setFilmOnePoster] = useState("");
  const [filmTwoPoster, setFilmTwoPoster] = useState("");
  const [filmThreePoster, setFilmThreePoster] = useState("");
  const [gameState, setGameState] = useState<gameState>();
  const [gameResponse, setgameResponse] = useState<guessResponse>();
  const [incorrectGuess, setIncorrectGuess] = useState<Set<string>>(new Set());
  const [correctGuess, setCorrectGuess] = useState<Set<string>>(new Set());

  const playerID = localStorage.getItem("playerID") ?? crypto.randomUUID();

  localStorage.setItem("playerID", playerID);

  async function gameStatePull() {
    try {
      const response = await fetch("http://localhost:8080/api/gamestate", {
        headers: {
          "Content-Type": "application/json",
          "X-Player-ID": playerID,
        },
      });
      const data = await response.json();

      for (const guess of data.guesses) {
        if (guess.Verdict === true) {
          switch (guess.FilmNumber) {
            case 1:
              setFilmOne(guess.Guess);
              setFilmOnePoster(`https://image.tmdb.org/t/p/w500${data.posters[0]}`);
              break;
            case 2:
              setFilmTwo(guess.Guess);
              setFilmTwoPoster(`https://image.tmdb.org/t/p/w500${data.posters[1]}`);
              break;
            case 3:
              setFilmThree(guess.Guess);
              setFilmThreePoster(`https://image.tmdb.org/t/p/w500${data.posters[2]}`);
              break;
          }
        }
      }

      setGameState(data);

      if (data) {
        for (const guess of data?.guesses) {
          guess.Verdict === true
            ? setCorrectGuess((prev) => new Set(prev).add(guess.Guess))
            : setIncorrectGuess((prev) => new Set(prev).add(guess.Guess));
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
    console.log(data);

    switch (data.film_number) {
      case 1:
        setFilmOne(data.guess);
        setFilmOnePoster(`https://image.tmdb.org/t/p/w500${data.poster_path[0]}`);
        break;
      case 2:
        setFilmTwo(data.guess);
        setFilmTwoPoster(`https://image.tmdb.org/t/p/w500${data.poster_path[0]}`);
        break;
      case 3:
        setFilmThree(data.guess);
        setFilmThreePoster(`https://image.tmdb.org/t/p/w500${data.poster_path[0]}`);
        break;
    }

    if (data.verdict == false) {
      setIncorrectGuess(new Set(incorrectGuess).add(guess));
    } else {
      setCorrectGuess(new Set(correctGuess).add(guess));
    }
    setgameResponse(data);

    if (response.ok) {
      setGuess("");
    }
  };

  useEffect(() => {
    gamePull();
    gameStatePull();
  }, []);

  console.log(gameState);

  return (
    <div>
      <div className="actor-name">{info?.actor}</div>
      <div className="poster-container">
        <div className="poster-card">
          <div className="poster">
            poster
            <img src={filmOnePoster} />
          </div>
          <div className="title">{filmOne ? filmOne : "???"}</div>
        </div>
        <div className="poster-card">
          <div className="poster">
            poster
            <img src={filmTwoPoster} />
          </div>
          <div className="title">{filmTwo ? filmTwo : "???"}</div>
        </div>
        <div className="poster-card">
          <div className="poster">
            poster
            <img src={filmThreePoster} />
          </div>
          <div className="title">{filmThree ? filmThree : "???"}</div>
        </div>
      </div>
      <div>
        {(gameResponse && gameResponse?.strikes >= 3) || (gameState?.strikes && gameState?.strikes >= 3) ? (
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
            <div>strikes: {gameResponse?.strikes ?? gameState?.strikes}</div>
          </div>
        )}
      </div>
    </div>
  );
}

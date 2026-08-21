import { useState } from "react";
import { useEffect } from "react";

type gameInfo = {
  actor: string;
  gamedate: string;
};

type gameResponse = {
  verdict: boolean;
  strikes: number;
  guess: string;
  playerid: string;
  repeat: boolean;
};

export function Home() {
  const [info, setInfo] = useState<gameInfo>();
  const [guess, setGuess] = useState("");
  const [gameResponse, setgameResponse] = useState<gameResponse>();

  let playerID = localStorage.getItem("playerID");

  if (!playerID) {
    playerID = crypto.randomUUID();
    localStorage.setItem("playerID", playerID);
  }

  // Need to pull game state

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
    setgameResponse(data);

    if (response.ok) {
      setGuess("");
    }

    console.log(data);
  };

  useEffect(() => {
    gamePull();
  }, []);

  return (
    <div>
      <div>{info?.actor}</div>
      <div>
        {gameResponse && gameResponse?.strikes >= 3 ? (
          "Game over"
        ) : (
          <div>
            <form onSubmit={submitGuess}>
              <label>Guess</label>
              <input type="guess" value={guess} onChange={(e) => setGuess(e.target.value)} placeholder="Film" />
            </form>
            <div>{gameResponse?.verdict === true ? `${gameResponse.guess}: correct` : `${gameResponse?.guess}: incorrect`}</div>
            <div>strikes: {gameResponse?.strikes}</div>
          </div>
        )}
      </div>
    </div>
  );
}

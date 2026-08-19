import { useState } from "react";
import { useEffect } from "react";

type gameInfo = {
  actor: string;
  gamedate: string;
};

export function Home() {
  const [info, setInfo] = useState<gameInfo>();
  const [guess, setGuess] = useState("");
  const [answer, setAnswer] = useState("");
  async function datarino() {
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
      },
      body: JSON.stringify({
        guess: guess,
        gamedate: info?.gamedate,
      }),
    });

    const data = await response.json();
    setAnswer(data);

    if (response.ok) {
      setGuess("");
    }
  };

  useEffect(() => {
    datarino();
  }, []);

  return (
    <div>
      <div>{info?.actor}</div>
      <form onSubmit={submitGuess}>
        <label>Guess</label>
        <input type="guess" value={guess} onChange={(e) => setGuess(e.target.value)} placeholder="Film" />
      </form>
      <div>{answer}</div>
    </div>
  );
}

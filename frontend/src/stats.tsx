import { useEffect, useState } from "react";

type stats = {
  games_played: number;
  perfect_wins: number;
  victories: number;
  zero_correct: number;
  one_correct: number;
  two_correct: number;
};

export function Stats() {
  const [stats, setStats] = useState<stats>();

  const playerID = localStorage.getItem("playerID") ?? crypto.randomUUID();

  localStorage.setItem("playerID", playerID);

  async function statsPull() {
    try {
      const response = await fetch("http://localhost:8080/api/stats", {
        headers: {
          "content-type": "application/json",
          "X-Player-ID": playerID,
        },
      });

      const data = await response.json();

      setStats(data);
    } catch (err) {
      console.error(err);
      throw err;
    }
  }

  useEffect(() => {
    statsPull();
  }, []);

  return (
    <div>
      <div>Stats</div>
      <div>Games Played: {stats?.games_played}</div>
      <div>Wins: {stats?.victories}</div>
      <div>Perfect Wins: {stats?.perfect_wins}</div>
      <div>1/3: {stats?.one_correct}</div>
      <div>2/3: {stats?.two_correct}</div>
      <div>0/3: {stats?.zero_correct}</div>
    </div>
  );
}

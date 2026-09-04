import { useEffect, useState } from "react";

type stats = {
  games_played: number;
  win_percentage: string;
  perfect_wins: number;
  victories: number;
  zero_correct: number;
  one_correct: number;
  two_correct: number;
};

type guessDistribution = {
  label: string;
  value: number;
};

export function Stats() {
  const [stats, setStats] = useState<stats>();
  const [statsDistribution, setStatsDistribution] = useState<guessDistribution[]>([]);

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
      setStatsDistribution([
        { label: "0/3", value: data.zero_correct },
        { label: "1/3", value: data.one_correct },
        { label: "2/3", value: data.two_correct },
        { label: "3/3", value: data.victories },
      ]);
    } catch (err) {
      console.error(err);
      throw err;
    }
  }
  const max = Math.max(...statsDistribution.map((stat) => stat.value));

  useEffect(() => {
    statsPull();
  }, []);

  return (
    <div className="stats-background">
      <div className="stats-title">Stats</div>
      <div className="all-stats">
        <div className="top-stats">
          <span>Games Played: {stats?.games_played}</span>
          <span>Win Percentage: {stats?.win_percentage}</span>
          <span>Perfect Games: {stats?.perfect_wins}</span>
        </div>
        <div className="distro">
          {statsDistribution.map((stat) => (
            <div className="distro-row" key={stat.label}>
              <span className="distro-label">{stat.label}</span>
              <div className="distro-track">
                <div className="distro-bar" style={{ width: max === 0 || max === undefined ? "0%" : `${(stat.value / max) * 100}%` }}>
                  {stat.value}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

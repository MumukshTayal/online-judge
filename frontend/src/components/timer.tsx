import React, { useEffect, useState } from 'react';

function Timer({ endTime }) {
  const calculateTimeLeft = () => {
    const difference = +new Date(endTime) - +new Date();
    let timeLeft = {};

    if (difference > 0) {
      timeLeft = {
        hours: Math.floor((difference / (1000 * 60 * 60)) % 24),
        minutes: Math.floor((difference / 1000 / 60) % 60),
        seconds: Math.floor((difference / 1000) % 60),
      };
    }

    return timeLeft;
  };

  const [timeLeft, setTimeLeft] = useState(calculateTimeLeft());

  useEffect(() => {
    const timer = setTimeout(() => {
      setTimeLeft(calculateTimeLeft());
    }, 1000);

    return () => clearTimeout(timer);
  });

  const formatTime = (value) => {
    return value < 10 ? `0${value}` : value;
  };

  // Check if time has passed
  const isContestEnded = +new Date(endTime) < +new Date();

  return (
    <div>
      {isContestEnded ? (
        <span>Contest ended</span>
      ) : (
        <>
          {timeLeft.hours > 0 && (
            <span>{`${formatTime(timeLeft.hours)}:`}</span>
          )}
          <span>{`${formatTime(timeLeft.minutes)}:`}</span>
          <span>{`${formatTime(timeLeft.seconds)}`}</span>
        </>
      )}
    </div>
  );
}

export default Timer;

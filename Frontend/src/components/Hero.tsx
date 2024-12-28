import React from "react";

const Hero: React.FC = () => {
  return (
    <section id="hero" className="hero">
      <div className="container">
        <h2>Explore the Waters with Ease</h2>
        <p>
          Rent boats, jet skis, and paddleboards for an unforgettable
          experience.
        </p>
        <a href="#services" className="btn">
          View Rentals
        </a>
      </div>
    </section>
  );
};

export default Hero;

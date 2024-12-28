import React from "react";
import Hero from "../components/Hero";
import About from "../components/About";
import Services from "../components/Services";

const HomePage: React.FC = () => {
  return (
    <>
      <Hero />
      <div className="home-wrapper">
        <About />
        <h2>Choose your spot</h2>
        <Services />
      </div>
    </>
  );
};

export default HomePage;

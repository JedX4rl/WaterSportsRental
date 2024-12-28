import React from "react";
import ServiceCard from "./ServiceCard";

interface Rental {
  image: string;
  title: string;
  description: string;
  slug: string;
}

const Services: React.FC = () => {
  const rentals: Rental[] = [
    {
      image: "boat.jpg",
      title: "Bali",
      description: "Relax in Bali's paradise.",
      slug: "bali",
    },
    {
      image: "jet-ski.jpg",
      title: "Hawaii",
      description: "Experience the thrill in Hawaii.",
      slug: "hawaii",
    },
    {
      image: "paddleboard.jpg",
      title: "Sri Lanka",
      description: "Explore calm waters in Sri Lanka.",
      slug: "sri-lanka",
    },
  ];

  return (
    <section id="services" className="services">
      {rentals.map((rental, index) => (
        <ServiceCard key={index} {...rental} />
      ))}
    </section>
  );
};

export default Services;

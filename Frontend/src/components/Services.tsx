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
      image: "src/assets/Bali.png",
      title: "Bali",
      description: "Relax in Bali's paradise.",
      slug: "1",
    },
    {
      image: "src/assets/Hawaii.png",
      title: "Hawaii",
      description: "Experience the thrill in Hawaii.",
      slug: "2",
    },
    {
      image: "src/assets/Sri-lanka.png",
      title: "Sri Lanka",
      description: "Explore calm waters in Sri Lanka.",
      slug: "3",
    },
  ];

  return (
    <body>
      <section id="services" className="services">
        {rentals.map((rental, index) => (
          <ServiceCard key={index} {...rental} />
        ))}
      </section>
    </body>
  );
};

export default Services;

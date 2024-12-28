import React from "react";
import { useNavigate } from "react-router";

interface ServiceCardProps {
  image: string;
  title: string;
  description: string;
  slug: string;
}

// const ServiceCard: React.FC<ServiceCardProps> = ({image, title, description}) => {
//     return (
//         <a href={""}>
//             <div className="card">
//                 <img src={image} alt={title}/>
//                 <h3>{title}</h3>
//                 <p>{description}</p>
//             </div>
//         </a>
//     );
// };
//
// export default ServiceCard;

const ServiceCard: React.FC<ServiceCardProps> = ({
  image,
  title,
  description,
  slug,
}) => {
  const navigate = useNavigate();

  // const handleClick = () => {
  //     navigate(`/products/?location=${slug}`); // Перенаправление на страницу
  // };
  const handleClick = () => {
    navigate(`/products?location=${slug}`); // Используем query string
  };

  return (
    <div className="card" onClick={handleClick} style={{ cursor: "pointer" }}>
      <img src={image} alt={title} />
      <h3>{title}</h3>
      <p>{description}</p>
    </div>
  );
};

export default ServiceCard;

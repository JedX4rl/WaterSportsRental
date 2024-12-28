// import React from "react";
// import {useNavigate} from "react-router";
//
// interface ProductCardProps {
//     image: string;
//     //     // title: string;
//     //     // description: string;
//     //     // slug: string;
//     id: number;
//     location_id: number;
//     type: string;
//     brand: string;
//     model: string;
//     year: number;
//     price: number;
// }
//
// const ProductCard: React.FC<ProductCardProps> = ({ image, location_id, type, brand, model, year, price }) => {
//     const navigate = useNavigate();
//     const handleClick = () => {
//         navigate(`/products?location=${location_id}`); // Используем query string
//     };
//
//     return (
//         <div className="card" onClick={handleClick} style={{ cursor: "pointer" }}>
//             <img src={image} alt={brand} />
//             <h3>{type}</h3>
//             <p>{model}</p>
//             <p>{year}</p>
//             <p>{price}</p>
//         </div>
//     );
// };
//
// export default ProductCard;

import React from "react";
import { useNavigate } from "react-router";

interface ProductCardProps {
  id: number;
  image: string;
  type: string;
  brand: string;
  model: string;
  year: number;
  price: number;
}

const ProductCard: React.FC<ProductCardProps> = ({
  id,
  image,
  type,
  brand,
  model,
  year,
  price,
}) => {
  const navigate = useNavigate();

  const handleClick = () => {
    navigate(`/products/${id}`);
  };

  return (
    <div className="product-card" onClick={handleClick}>
      <div className="card-image">
        <img src={image} alt={`${brand} ${model}`} />
      </div>
      <div className="card-content">
        <h3>
          {brand} {model}
        </h3>
        <p>Type: {type}</p>
        <p>Year: {year}</p>
        <p>
          Price: <strong>${price}/hour</strong>
        </p>
      </div>
    </div>
  );
};

export default ProductCard;

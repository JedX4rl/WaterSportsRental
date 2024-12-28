import React from "react";
import ProductCard from "./ProductCard";

interface Rental {
  id: number;
  location_id: number;
  type: string;
  brand: string;
  model: string;
  year: number;
  price: number;
}

interface ProductsProps {
  rentals: Rental[];
}

const Products: React.FC<ProductsProps> = ({ rentals }) => {
  return (
    <section id="products" className="products">
      {rentals?.length === 0 ? (
        <div>No products available</div> // Если нет товаров
      ) : (
        rentals.map((rental) => (
          <ProductCard image={""} key={rental.id} {...rental} />
        ))
      )}
    </section>
  );
};

export default Products;

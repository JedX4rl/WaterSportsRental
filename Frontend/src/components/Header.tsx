import React from "react";
import { Link } from "react-router";

const Header: React.FC = () => {
  return (
    <header className="header">
      <div className="logo">
        <Link to="/">WaterSports Rental</Link>
      </div>
      <nav className="nav">
        <ul>
          <li>
            <Link to="/">Home</Link>
          </li>
          <li>
            <a href="/#about">About</a>
          </li>
          <li>
            <a href="/#services">Services</a>
          </li>
          <li>
            <Link to={"/contacts"}>Contacts</Link>
          </li>
          <li>
            <Link to="/products">Products</Link>
          </li>
          <li>
            <Link to="/auth">Profile</Link>
          </li>
        </ul>
      </nav>
    </header>
  );
};

export default Header;

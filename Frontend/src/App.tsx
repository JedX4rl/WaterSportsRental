import React from "react";
import { Routes, Route } from "react-router";
import Header from "./components/Header";
import Footer from "./components/Footer";
import HomePage from "@/pages/HomePage";
import ProductsPage from "@/pages/ProductsPage";
import ProductDetailsPage from "@/components/ProductDetailsPage.tsx";
import UserProfile from "@/components/UserProfile.tsx";
import Auth from "@/components/Auth.tsx";
import AdminPage from "@/components/AdminPage.tsx";
import ItemsPage from "@/components/ManageItems.tsx";
import LocationsPage from "@/components/ManageLocations.tsx";
import UsersPage from "@/components/ManageUsers.tsx";
import ReviewsPage from "@/components/ManageReviews.tsx";
import "./index.css"
import DumpsPage from "@/components/ManageDumps.tsx";
import Contacts from "@/components/Contacts.tsx";

const App: React.FC = () => {
  return (
    <React.Fragment>
      <Header />
      <main>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="products" element={<ProductsPage />} />{" "}
          <Route path="products/:id" element={<ProductDetailsPage />} />
          <Route path="profile" element={<UserProfile />} />
          <Route path="auth" element={<Auth />} />
          <Route path="admin" element={<AdminPage />} />
          <Route path="*" element={<h1>404 - Page Not Found</h1>} />
          <Route path="admin/items" element={<ItemsPage />} />
          <Route path="admin/locations" element={<LocationsPage />} />
          <Route path="admin/users" element={<UsersPage />} />
          <Route path="admin/reviews" element={<ReviewsPage />} />
          <Route path="admin/dumps" element={<DumpsPage />} />
          <Route path="contacts" element={<Contacts/>} />
        </Routes>
      </main>
      <Footer />
    </React.Fragment>
  );
};

export default App;

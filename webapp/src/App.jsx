import React from "react";
import "./App.css";
import { useState } from "react";
import ModalDialog from "./components/ModalDialog";
import { createBrowserRouter } from "react-router-dom";

import { RouterProvider } from "react-router-dom";

import Home from "./routes/Home";
import NewPoll from "./routes/NewPoll";
import Form from "./components/SignupForm";
import Signup from "./routes/SignUp";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Home />,
  },
  {
    path: "/new",
    element: <NewPoll />,
  },
  {
    path: "/signup",
    element: <Signup />,
  },
]);

export function App() {
  return (
    <div className="app">
      <React.StrictMode>
        <RouterProvider router={router} />
      </React.StrictMode>
    </div>
  );
}

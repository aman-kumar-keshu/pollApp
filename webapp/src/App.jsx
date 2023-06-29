import React from "react";
import "./App.css";
import { createBrowserRouter } from "react-router-dom";

import { RouterProvider } from "react-router-dom";

import Home from "./routes/Home";
import NewPoll from "./routes/NewPoll";
import Signup from "./routes/SignUp";
import Login from "./routes/Login";

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
  {
    path: "/login",
    element: <Login />,
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

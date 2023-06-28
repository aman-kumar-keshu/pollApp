import React from "react";

import { useNavigate } from "react-router-dom";
import Button from "@material-ui/core/Button";
const Header = (props) => {
  const navigate = useNavigate();

  return (
    <>
      <div
        style={{
          position: "sticky",
          top: 0,
          padding: "10px 30px",
          height: "40px",
          background: "black",
          color: "white",
        }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
          }}
        >
          <>
            <div>polling-app</div>{" "}
            {token ? (
              <Button variant="contained" color="primary" onClick={handleOpen}>
                Logout
              </Button>
            ) : (
              <Button variant="contained" color="primary" onClick={handleOpen}>
                Signup
              </Button>
            )}
            {token && (
              <button
                onClick={() => {
                  navigate("/new");
                }}
                style={{ height: "50px" }}
              >
                Create New poll
              </button>
            )}
          </>
        </div>
      </div>
      <div className="card-header"></div>
      <h1
        style={{
          marginLeft: "500px",
        }}
        className="card-header-title header"
      >
        You have {props.numPolls} Polls
      </h1>
    </>
  );
};
export default Header;

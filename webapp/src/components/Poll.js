import React from "react";

function Poll(props) {
  return (
    <div
      className="content"
      style={{
        margin: "auto",
      }}
    >
      <li>
        <div>
          <h2>{props.poll.name} </h2>
        </div>
        <div>
          {" "}
          <h4> topic : {props.poll.topic}</h4>
        </div>
        <div>
          {" "}
          <img
            height={"100px"}
            width={"100px"}
            src={
              props.poll.src
              // "https://i.natgeofe.com/n/e76f5368-6797-4794-b7f6-8d757c79ea5c/ng-logo-2fl.png?w=109&h=32"
            }
          />
        </div>
        <button
          onClick={() => {
            props.onUpdate(props.poll.id, {
              ...props.poll,
              upvotes: props.poll.upvotes + 1,
            });
          }}
        >
          <strong>upvotes: </strong> {props.poll.upvotes}{" "}
        </button>

        <button
          onClick={() => {
            props.onUpdate(props.poll.id, {
              ...props.poll,
              downvotes: props.poll.downvotes + 1,
            });
          }}
        >
          {" "}
          <strong>downvotes: </strong> {props.poll.downvotes}{" "}
        </button>
      </li>

      <button
        onClick={() => {
          props.onDelete(props.poll.id);
        }}
      >
        <strong> Delete Poll</strong>
      </button>
    </div>
  );
}

export default Poll;

import React from "react";

function Poll(props) {
  return (
    <div className="list-item">
      <li>
        <div>Name : {props.poll.name}</div>
        <div>topic : {props.poll.topic}</div>
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
          upvotes : {props.poll.upvotes}{" "}
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
          downvotes : {props.poll.downvotes}{" "}
        </button>
      </li>

      <button
        onClick={() => {
          props.onDelete(props.poll.id);
        }}
      >
        {" "}
        Delete Poll
      </button>
    </div>
  );
}

export default Poll;

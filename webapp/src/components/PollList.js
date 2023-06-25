import React, { Component } from "react";
import Poll from "./Poll";
function PollList(props) {
  return (
    <div className="list-wrapper">
      {props.polls.map((poll) => (
        <Poll
          poll={poll}
          key={poll.id}
          id={poll.id}
          onUpdate={props.onUpdate}
          onDelete={props.onDelete}
        />
      ))}
    </div>
  );
}
export default PollList;

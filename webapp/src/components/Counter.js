import { Component } from "react";

class Counter extends Component {
  state = {
    count: 0,
  };

  handleClick = () => {
    this.setState(({ count }) => ({
      count: count + 1,
    }));
    


  };

  render() {
    return (
      <div>
        <button onClick={this.handleClick}>{this.state.count}</button>
      </div>
    );
  }
}
export default Counter;

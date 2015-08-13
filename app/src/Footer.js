import React, { Component } from 'react';

export default class App extends Component {
  render() {
    var footerStyle = {
      bottom: 20,
      float: 'right'
    };
    return (
      <div style={footerStyle}>
        <a className="ui large label">
          <i className="twitter icon"></i> @crosbymichael
        </a>
        <a className="ui large label">
          <i className="twitter icon"></i> @arunan
        </a>
      </div>
    );
  }
}

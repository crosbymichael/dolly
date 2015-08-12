import React, { Component } from 'react';
import ContainersGrid from './ContainersGrid';
import TopNav from './TopNav';

export default class App extends Component {
  render() {
    return (
      <div>
        <TopNav />
        <ContainersGrid />
      </div>
    );
  }
}

import React, { Component } from 'react';
import ContainerItem from './ContainerItem';

export default class ContainersGrid extends Component {
  render() {
    return (
      <div className="ui cards six column padded grid">
        <ContainerItem />
        <ContainerItem />
        <ContainerItem />
        <ContainerItem />
        <ContainerItem />
        <ContainerItem />
        <ContainerItem />
        <ContainerItem />
      </div>
    );
  }
}

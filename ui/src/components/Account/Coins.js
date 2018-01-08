import React, { Component } from "react"

import { Header, Segment, Label } from "semantic-ui-react"

class Coins extends Component {
  render() {
    const title = this.props.title
    const coins = this.props.coins

    return (
      <div>
        <Header>{title}</Header>
        {coins.map((coin, index) => (
          <Segment key={index} compact vertical>
            <Label>
              {coin.amount} <Label.Detail>{coin.denom}</Label.Detail>
            </Label>
          </Segment>
        ))}
      </div>
    )
  }
}

export default Coins

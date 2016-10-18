import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';

import New from './New';

@connect(createStructuredSelector({
  gamesOptions: selectors.gamesOptions,
}), actions)
export default class GameMapNew extends Component {
  static propTypes = {
    gamesOptions: ImmutablePropTypes.list,

    loadGames: PropTypes.func.isRequired,
    createGameMap: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.props.loadGames();
  }

  shouldComponentUpdate(nextProps) {
    return nextProps.gamesOptions !== this.props.gamesOptions;
  }

  handleSuccess = () => {
    // TODO: browserHistory.push(`/game-maps/${x.id}`);
    this.props.messagePush('created a game map');
  };

  render() {
    return (
      <New
        isAdminOnly
        resourceName="game map"
        fields={{
          gameId: {
            name: 'the game',
            type: 'select',
            options: this.props.gamesOptions.toJS(),
          },
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          sideX: { name: 'Side 1' },
          sideXAbbr: { name: 'Side 1 abbreviation' },
          sideY: { name: 'Side 2' },
          sideYAbbr: { name: 'Side 2 abbreviation' },
        }}
        createResource={this.props.createGameMap}
        onSuccess={this.handleSuccess}
      />
    );
  }
}

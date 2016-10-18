import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';

import { gameMapsOptionsSelector } from './MatchPrepare';

import New from './New';

@connect(createStructuredSelector({
  gameMapsOptions: gameMapsOptionsSelector,
}), actions)
export default class BracketMapNew extends Component {
  static propTypes = {
    game: ImmutablePropTypes.map,
    bracket: ImmutablePropTypes.map,

    gameMapsOptions: ImmutablePropTypes.list,

    loadGameMaps: PropTypes.func.isRequired,
    createBracketMap: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ game, gameMapsOptions }, prevProps) {
    if (!game) return;
    if (
      !prevProps ||
      game !== prevProps.game ||
      (!gameMapsOptions && prevProps.gameMapsOptions)
    ) {
      this.props.loadGameMaps({ gameId: game.get('id') });
    }
  }

  render() {
    const { bracket, gameMapsOptions } = this.props;
    if (!bracket || !gameMapsOptions) return null;
    return (
      <New
        isAdminOnly
        resourceName="bracket map"
        fields={{
          gameMapId: {
            name: 'The map',
            type: 'select',
            options: gameMapsOptions.toJS(),
          },
          subPool: { name: 'Sub-pool 0-9', value: '0' },
          bracketId: { type: 'hidden', value: bracket.get('id') },
        }}
        createResource={this.props.createBracketMap}
      />
    );
  }
}

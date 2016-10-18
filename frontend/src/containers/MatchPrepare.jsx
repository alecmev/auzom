import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';

import New from './New';

export const gameMapsOptionsSelector = selectors.options(createSelector(
  selectors.gameMaps,
  (_, props) => props.game && props.game.get('id'),
  (x, gameId) => gameId && x.toList().filter(y => y.get('gameId') === gameId),
));

@connect(createStructuredSelector({
  gameMapsOptions: gameMapsOptionsSelector,
}), actions)
export default class MatchPrepare extends Component {
  static propTypes = {
    game: ImmutablePropTypes.map,
    matchPath: PropTypes.string,
    match: ImmutablePropTypes.map,

    gameMapsOptions: ImmutablePropTypes.list,

    loadGameMaps: PropTypes.func.isRequired,
    patchMatch: PropTypes.func.isRequired,
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

  patcher = (...args) => {
    if (!this.props.match) return;
    this.props.patchMatch(this.props.match.get('id'), ...args);
  };

  // intentional force reload, because easy
  handleSuccess = () => { location.href = this.props.matchPath; };

  render() {
    const { gameMapsOptions } = this.props;
    return (
      <New
        lighter
        isAdminOnly
        resourceName="match"
        fields={{
          maps: {
            name: 'The maps',
            type: 'multi',
            options: gameMapsOptions ? gameMapsOptions.toJS() : [],
          },
          action: { type: 'hidden', value: 'prepare' },
        }}
        callForAction="prepare the match"
        createResource={this.patcher}
        onSuccess={this.handleSuccess}
      />
    );
  }
}

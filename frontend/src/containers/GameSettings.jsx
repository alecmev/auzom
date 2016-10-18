import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';

import Settings from './Settings';

const gameSlugSelector = (state, props) => props.params.slug;
const gameSelector = createSelector(
  selectors.games,
  gameSlugSelector,
  (games, gameSlug) =>
    games.toList().find(x => x.get('slug') === gameSlug),
);

@connect(createStructuredSelector({
  gameSlug: gameSlugSelector,
  game: gameSelector,
}), actions)
export default class GameSettings extends Component {
  static propTypes = {
    gameSlug: PropTypes.string.isRequired,
    game: ImmutablePropTypes.map,

    loadGames: PropTypes.func.isRequired,
    updateGame: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ gameSlug, game }, prevProps) {
    if (
      !prevProps ||
      gameSlug !== prevProps.gameSlug ||
      (!game && prevProps.game)
    ) {
      this.props.loadGames({ slug: gameSlug });
    }
  }

  handleSuccess = (m, x) => {
    m.slug && browserHistory.push(`/${x.slug}/settings`);
    this.props.messagePush('game settings saved');
  };

  render() {
    return (
      <Settings
        isAdminOnly
        resourceName="game"
        resource={this.props.game}
        fields={{
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          slug: { name: 'URL-friendly name' },
          cover: { name: 'Cover art URL' },
          summary: { name: 'Summary', type: 'textarea' },
          releasedAt: { type: 'date', name: 'Release date' },
        }}
        updateResource={this.props.updateGame}
        onSuccess={this.handleSuccess}
      />
    );
  }
}

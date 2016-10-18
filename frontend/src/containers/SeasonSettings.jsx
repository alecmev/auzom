import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';
import { createSelector, createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';

import Settings from './Settings';

const seasonSlugSelector = (state, props) => props.params.seasonSlug;
const seasonSelector = createSelector(
  selectors.seasons,
  (state, props) => props.tournament,
  seasonSlugSelector,
  (seasons, tournament, slug) => seasons.toList().find(x =>
    x.get('tournamentId') === tournament.get('id') &&
    x.get('slug') === slug
  ),
);

@connect(createStructuredSelector({
  seasonSlug: seasonSlugSelector,
  season: seasonSelector,
}), actions)
export default class SeasonSettings extends Component {
  static propTypes = {
    tournamentPath: PropTypes.string.isRequired,
    tournament: ImmutablePropTypes.map,
    season: ImmutablePropTypes.map,

    loadSeasons: PropTypes.func.isRequired,
    updateSeason: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.load(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.load(nextProps, this.props);
  }

  load({ tournament, seasonSlug, season }, prevProps) {
    if (!tournament) return;
    if (
      !prevProps ||
      tournament !== prevProps.tournament ||
      seasonSlug !== prevProps.seasonSlug ||
      (!season && prevProps.season)
    ) {
      this.props.loadSeasons({
        tournamentId: tournament.get('id'),
        slug: seasonSlug,
      });
    }
  }

  handleSuccess = (m, x) => {
    m.slug && browserHistory.push(
      `${this.props.tournamentPath}/${x.slug}/_/settings`,
    );
    this.props.messagePush('season settings saved');
  };

  render() {
    return (
      <Settings
        lighter
        isAdminOnly
        resourceName="season"
        resource={this.props.season}
        fields={{
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          slug: { name: 'URL-friendly name' },
          publishedAt: { name: 'Publish time', type: 'datetime' },
          description: { name: 'One-paragraph description', type: 'textarea' },
          teamSize: { name: 'Min team size / format size', type: 'number' },
          teamSizeMax: { name: 'Max team size', type: 'number' },
          capacity: { name: 'Max teams', type: 'number' },
          duration: { name: 'Approximate duration, days', type: 'number' },
          youtubePlaylist: { name: 'YouTube playlist' },
          sponsors: {
            name: 'Sponsors [title;logo;scale;url]',
            type: 'textarea',
          },
          signupsOpenedAt: { type: 'datetime', name: 'Signups open time' },
          signupsClosedAt: { type: 'datetime', name: 'Signups close time' },
          endedAt: { type: 'datetime', name: 'Season end time' },
        }}
        updateResource={this.props.updateSeason}
        onSuccess={this.handleSuccess}
      />
    );
  }
}

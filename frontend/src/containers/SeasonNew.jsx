import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import * as actions from '../actions';
import * as selectors from '../selectors';

import New from './New';

@connect(createStructuredSelector({
  tournamentsOptions: selectors.tournamentsOptions,
}), actions)
export default class SeasonNew extends Component {
  static propTypes = {
    tournamentsOptions: ImmutablePropTypes.list,

    loadTournaments: PropTypes.func.isRequired,
    createSeason: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  componentWillMount() {
    this.props.loadTournaments();
  }

  shouldComponentUpdate(nextProps) {
    return nextProps.tournamentsOptions !== this.props.tournamentsOptions;
  }

  handleSuccess = () => {
    // TODO: browserHistory.push(`/seasons/${x.id}`);
    this.props.messagePush('created a season');
  };

  render() {
    return (
      <New
        isAdminOnly
        resourceName="season"
        fields={{
          tournamentId: {
            name: 'The tournament',
            type: 'select',
            options: this.props.tournamentsOptions.toJS(),
          },
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          slug: { name: 'URL-friendly name' },
          publishedAt: { name: 'Publish time', type: 'datetime' },
          description: { name: 'One-paragraph description', type: 'textarea' },
          rules: { name: 'Rules', type: 'textarea' },
          teamSize: { name: 'Min team size / format size' },
          teamSizeMax: { name: 'Max team size' },
          capacity: { name: 'Max teams' },
          duration: { name: 'Approximate duration, days' },
          youtubePlaylist: { name: 'YouTube playlist ID' },
          sponsors: {
            name: 'Sponsors [title;logo;scale;url]',
            type: 'textarea',
          },
          signupsOpenedAt: { name: 'Signups open time', type: 'datetime' },
          signupsClosedAt: { name: 'Signups close time', type: 'datetime' },
          endedAt: { name: 'Season end time', type: 'datetime' },
        }}
        createResource={this.props.createSeason}
        onSuccess={this.handleSuccess}
      />
    );
  }
}

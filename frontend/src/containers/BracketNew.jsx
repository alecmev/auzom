import { Component, PropTypes } from 'react';
import ImmutablePropTypes from 'react-immutable-proptypes';
import { connect } from 'react-redux';

import * as actions from '../actions';

import New from './New';

@connect(undefined, actions)
export default class BracketNew extends Component {
  static propTypes = {
    stagePath: PropTypes.string,
    stage: ImmutablePropTypes.map,

    createBracket: PropTypes.func.isRequired,
    messagePush: PropTypes.func.isRequired,
  };

  static contextTypes = {
    router: PropTypes.object.isRequired,
  };

  handleSuccess = (x) => {
    this.context.router.push(`${this.props.stagePath}/${x.slug}`);
    this.props.messagePush('created a bracket');
  };

  render() {
    if (!this.props.stage) return null;
    return (
      <New
        isAdminOnly
        resourceName="bracket"
        fields={{
          name: { name: 'Name' },
          abbr: { name: 'Abbreviation' },
          slug: { name: 'URL-friendly name' },
          order: { name: 'Sort number' },
          type: {
            name: 'Bracket type',
            type: 'select',
            options: [
              { value: 'bcl-s8-group-stage', label: 'Round robin' },
              { value: 'bcl-s8-playoffs', label: 'Single elimination' },
              { value: 'bcl-sc16-swiss', label: 'BCL Swiss' },
              { value: 'ace-pre-swiss', label: 'ACE Swiss' },
            ],
          },
          size: { name: 'Number of participants' },
          mapVetoProcedure: { name: 'Map veto procedure', type: 'textarea' },
          startAt: { name: 'First match time', type: 'datetime' },
          waitDays: { name: 'Days between rounds' },
          sameDayWaitMinutes: { name: 'Minutes between rounds on same day' },
          reportMinutes: { name: 'Minutes for report submission' },
          stageId: { type: 'hidden', value: this.props.stage.get('id') },
        }}
        createResource={this.props.createBracket}
        onSuccess={this.handleSuccess}
      />
    );
  }
}

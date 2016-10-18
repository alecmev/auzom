import Component from '../utils/Component';

import Markdown from '../components/Markdown';

const markdown = `
**What is Auzom?** Auzom is a community driven esports organisation.

**The brand:** Auzom is the proto-Italic noun for gold. A suitable name for an
esports organization where players compete for first place. When pronounced it
sounds similar to awesome!

**What differentiates Auzom from the others?** Community & Quality. Community
because Auzom’s events are organized by competent players from within the
community. Quality because the Auzom events are long term projects supervised by
professionals, hosted and broadcasted with state of the art software.

**Auzom’s goal** is to expand the number of community supported leagues it hosts
into multiple games, game platforms and regions. Bringing together players and
teams.

**Who is behind Auzom?** The development team behind Auzom is also the team that
created and ran LevelBF, a well known Battlefield community that brought
extensive competitive coverage of it's events since the start of Battlefield 3
through in-house developed spectator tools. With more ambition in esports they
decided to rebrand the organization and continue as Auzom.

*The Team members are:*

 - Alexandre ‘UneFriteUneFois’ Niset - BCL Head Admin / Community Manager
 - Christiaan ‘Aether’ Arnaud - CEO / Co-founder
 - Greer ‘Daskro’ Carper - Board Advisor / Co-founder
 - Kevin ‘Kevinario’ van der Staaij - CPO / Co-founder
 - Olegs ‘jevs’ Jeremejevs - CTO / Co-founder

**Current gaming communities hosted by Auzom:**

 - *Battlefield Conquest League,* BCL is the largest Battlefield 8v8 tournament
   in Europe. Formerly known as Battlefield Nordic League and dating back to
   Battlefield 2. Recently BCL became available for the America's.
 
 - *Aces High,* A Battlefield PC community organising cups and leagues for
   Battlefield 3, Battlefield 4, and the soon to be released Battlefield 1.

**This website** is in development. Auzom has big plans for it and new features
will be released as work progresses.

**Feel free to contact us** if you have any inquiries or want to contribute in
any way: [support@auzom.gg](mailto:support@auzom.gg)
`;

export default class About extends Component {
  render() {
    return (
      <div className={this.cni({ u: 'sectionSingle' })}>
        <div className={this.cn({ u: 'sectionMargined' })}>
          <div className={this.cn({ d: 'inner' })}>
            <Markdown source={markdown} />
          </div>
        </div>
      </div>
    );
  }
}

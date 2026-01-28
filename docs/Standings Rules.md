# Standings and Tiebreaker Rules

## NFL Rules

### Regular Season Standings
* Teams are grouped by division (4 teams per division) and conference (4 divisions per conference)
* Win percentage = Wins / (Wins + Losses + 0.5 * Ties)
* Ties count as one-half win and one-half loss for both teams

### Playoff Seeding
* 7 teams per conference make the playoffs
* **Seeds 1-4**: Division winners (ranked by record)
* **Seeds 5-7**: Wild card teams (ranked by record)
* **Seed 1**: Gets first-round bye

### Tiebreaker Procedures

* Only one team advances in any tiebreaking step. Remaining tied teams revert back to first step of applicable procedure.
* In comparing records against common opponents amongst tied teams, best win percentage is deciding factor because the teams may have played an unequal number of games against common opponents.
* To determine tiebreakers among division winners, apply inter-conference tiebreakers.
* To determine tiebreakers among wild card teams, apply division tiebreakers if the involved teams are from the same division or conference tiebreakers otherwise.
* Treats a 0-0 record as less than 0.000. So, a team that is 0-0 will be lower than a team that is 0-1. This is particularly important early in the season when few games have been played.
* A game must have an outcome for it to count towards record against common opponents, as well as a team's strength of schedule and strength of victory.

#### Inter-Division Tiebreakers
* If two teams in the same division have an identical win percentage, these are the tiebreaking steps to be taken until a winner is determined:
    1. Head-to-head record
    2. Division record
    3. Record vs. common opponents
    4. Conference record
    5. Strength of victory
    6. Strength of schedule
    7. Point differential
    8. Points scored
    9. Points allowed
    10. Coin toss (random choice)
* If three or more teams in the same division have an identical win percentage, these are the tiebreaking steps to be taken until a winner is determined:
    * Note: If at some point during these steps, at least one team is eliminated and there are only two teams left that remain tied, the tiebreaker restarts at Step 1 of the two-team format above. If at some point during these steps, at least one team is eliminated but there are still at least three teams left that remain tied, the tiebreaker restarts at Step 1 of these steps.
    1. Head-to-head record
    2. Division record
    3. Record vs. common opponents
    4. Conference record
    5. Strength of victory
    6. Strength of schedule
    7. Coin toss (random choice)

#### Inter-Conference Tiebreakers
* If two teams in the same conference have an identical win percentage, these are the tiebreaking steps to be taken until a winner is determined:
    * Note: If the tied teams are from the same division, apply the inter-division tiebreaker. Otherwise, continue with these steps.
    1. Head-to-head record
    2. Conference record
    3. Record vs. common opponents (minimum 4 games)
    4. Strength of victory
    5. Strength of schedule
    6. Point differential
    7. Points scored
    8. Points allowed
    9. Coin toss (random choice)
* If three or more teams in the same division have an identical win percentage, these are the tiebreaking steps to be taken until a winner is determined:
    * Note: Note: If at some point during these steps, at least one team is eliminated and there are only two teams left that remain tied, the tiebreaker restarts at Step 1 of the two-team format above. If at some point during these steps, at least one team is eliminated but there are still at least three teams left that remain tied, the tiebreaker restarts at Step 1 of these steps.
    1. Apply inter-division tiebreaker to eliminate all but the highest ranked team in each division involved in tiebreaker. The original seeding within the division upon application of the division tiebreaker remains the same for all subsequent applications of the procedure that are necessary.
    2. Head-to-head sweep (applicable only if one team has defeated each of the others or if one team has lost to each of the others)
    3. Conference record
    4. Record vs. common opponents (minimum 4 games)
    5. Strength of victory
    6. Strength of schedule
    7. Coin toss (random choice)

### Algorithm

Pseudocode for standings calculation:

function calculateNFLStandings(scenario):
1. Get all games for season
2. Merge actual results of games with user picks (user picks take priority)
3. Calculate necessary statistics for each team
4. Group teams by division and conference
5. Rank division winners per conference
6. Rank wild cards
7. Apply tiebreakers
8. Return seeded teams

---

## NBA Rules

### Regular Season Standings
* Teams are grouped by division (5 teams per division) and conference (3 divisions per conference)
* Win percentage = Wins / (Wins + Losses + 0.5 * Ties)
* Ties count as one-half win and one-half loss for both teams

### Playoff Seeding
* 10 teams per conference make the postseason
* **Seeds 1-6**: Automatic playoff berths
* **Seeds 7-10**: Participants of single-elimination Play-In Tournament
* Playoff series are best-of-7

### Tiebreaker Procedures

* Ties to determine division winners must be broken before any other ties.
* When a tie must be broken to determine a division winner, the results of the tie-break shall be used to determine only the division winner, and not for any other purpose.
* For multi-team tiebreaking process:
    * **Complete Breaking**: If each tied team has a different win percentage or point differential under the applicable criterion, teams are ranked accordingly and no further tiebreakers are needed.
    * **Partial Breaking**: If one or more (but not all) teams have different performance under a criterion:
        * Better performing team(s) get higher playoff position(s)
        * Remaining tied teams restart the tiebreaker process from the beginning using two-team criteria (if two teams remain) or multi-team criteria (if three or more remain)
    * **Random Drawing**: If application of all criteria does not break the tie, playoff positions are determined by random drawing.

#### Two-Way Tiebreakers
* In the case of a tie in regular season records involving only two teams, the following criteria will be utilized in the following order:
    1. Head-to-head record
    2. Division winner (this criterion is applied regardless of whether the tied teams are in the same division)
    3. Division record (only if tied teams are in same division)
    4. Conference record
    5. Win percentage vs. teams eligible for postseason in own conference
    6. Win percentage vs. teams eligible for postseason in opposing conference
    7. Point differential

#### Multi-Way Tiebreakers
* In the case of a tie in regular season records involving more than two teams, the following criteria will be utilized in the following order:
    1. Division record (this criterion is applied regardless of whether the tied teams are in the same division)
    2. Head-to-head record among tied teams
    3. Division record (only if all tied teams are in same division)
    4. Conference record
    5. Win percentage vs. teams eligible for postseason in own conference
    6. Point differential

---

## CFB Rules

### Regular Season Rankings

### College Football Playoff

### Conference Standings

### Conference Tiebreakers

### Bowl Game Selection

---


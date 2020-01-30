Feature: General

  Scenario: Create banner1
    When I request to create banner 1 with data:
      """
      {
          "id": 1,
          "slot_id": 1,
          "note": "banner note"
      }
      """
    Then The response code should be 200

  Scenario: Create banner2
    When I request to create banner 2 with data:
      """
      {
          "id": 2,
          "slot_id": 1,
          "note": "banner note"
      }
      """
    Then The response code should be 200

  Scenario: Create banner3
    When I request to create banner 3 with data:
      """
      {
          "id": 3,
          "slot_id": 1,
          "note": "banner note"
      }
      """
    Then The response code should be 200

  Scenario: Select N times
    When I request banner to show for 90 times
    Then I will get 30 shows for banner 1
    And I will get 30 shows for banner 2
    And I will get 30 shows for banner 3

  Scenario: Send click
    When I request click for banner 1:
    Then The response code should be 200

  Scenario: Send click
    When I request click for banner 1:
    Then The response code should be 200

  Scenario: Send click
    When I request click for banner 2:
    Then The response code should be 200

  Scenario: Select N times
    When I request banner to show for 90 times
    Then banner 1 show count must be greater than banner 2 show count
    And banner 2 show count must be greater than banner 3 show count

  Scenario: Delete banner1
    When I request to delete banner 1
    Then The response code should be 200

# coding: utf-8

"""
    Hatchet API

    The Hatchet API

    The version of the OpenAPI document: 1.0.0
    Generated by OpenAPI Generator (https://openapi-generator.tech)

    Do not edit the class manually.
"""  # noqa: E501


from __future__ import annotations
import pprint
import re  # noqa: F401
import json

from pydantic import BaseModel, Field, StrictInt, StrictStr
from typing import Any, ClassVar, Dict, List
from hatchet_sdk.clients.rest.models.pull_request_state import PullRequestState
from typing import Optional, Set
from typing_extensions import Self

class PullRequest(BaseModel):
    """
    PullRequest
    """ # noqa: E501
    repository_owner: StrictStr = Field(alias="repositoryOwner")
    repository_name: StrictStr = Field(alias="repositoryName")
    pull_request_id: StrictInt = Field(alias="pullRequestID")
    pull_request_title: StrictStr = Field(alias="pullRequestTitle")
    pull_request_number: StrictInt = Field(alias="pullRequestNumber")
    pull_request_head_branch: StrictStr = Field(alias="pullRequestHeadBranch")
    pull_request_base_branch: StrictStr = Field(alias="pullRequestBaseBranch")
    pull_request_state: PullRequestState = Field(alias="pullRequestState")
    __properties: ClassVar[List[str]] = ["repositoryOwner", "repositoryName", "pullRequestID", "pullRequestTitle", "pullRequestNumber", "pullRequestHeadBranch", "pullRequestBaseBranch", "pullRequestState"]

    model_config = {
        "populate_by_name": True,
        "validate_assignment": True,
        "protected_namespaces": (),
    }


    def to_str(self) -> str:
        """Returns the string representation of the model using alias"""
        return pprint.pformat(self.model_dump(by_alias=True))

    def to_json(self) -> str:
        """Returns the JSON representation of the model using alias"""
        # TODO: pydantic v2: use .model_dump_json(by_alias=True, exclude_unset=True) instead
        return json.dumps(self.to_dict())

    @classmethod
    def from_json(cls, json_str: str) -> Optional[Self]:
        """Create an instance of PullRequest from a JSON string"""
        return cls.from_dict(json.loads(json_str))

    def to_dict(self) -> Dict[str, Any]:
        """Return the dictionary representation of the model using alias.

        This has the following differences from calling pydantic's
        `self.model_dump(by_alias=True)`:

        * `None` is only added to the output dict for nullable fields that
          were set at model initialization. Other fields with value `None`
          are ignored.
        """
        excluded_fields: Set[str] = set([
        ])

        _dict = self.model_dump(
            by_alias=True,
            exclude=excluded_fields,
            exclude_none=True,
        )
        return _dict

    @classmethod
    def from_dict(cls, obj: Optional[Dict[str, Any]]) -> Optional[Self]:
        """Create an instance of PullRequest from a dict"""
        if obj is None:
            return None

        if not isinstance(obj, dict):
            return cls.model_validate(obj)

        _obj = cls.model_validate({
            "repositoryOwner": obj.get("repositoryOwner"),
            "repositoryName": obj.get("repositoryName"),
            "pullRequestID": obj.get("pullRequestID"),
            "pullRequestTitle": obj.get("pullRequestTitle"),
            "pullRequestNumber": obj.get("pullRequestNumber"),
            "pullRequestHeadBranch": obj.get("pullRequestHeadBranch"),
            "pullRequestBaseBranch": obj.get("pullRequestBaseBranch"),
            "pullRequestState": obj.get("pullRequestState")
        })
        return _obj


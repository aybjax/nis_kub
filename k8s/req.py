import requests
import random
import string
from concurrent.futures import ThreadPoolExecutor

def randomString():
    length = random.randint(3, 20)
    return ''.join(random.choice(string.ascii_letters ) for x in range(length))

student_ids = {}
course_ids = {}
student_arr = []
course_arr = []

students_api = 'http://localhost:8000/api/students'
courses_api = 'http://localhost:8000/api/courses'

def studentsCreate():
    resp = requests.post(
        f'{students_api}',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
        json={
            "name": randomString(),
            "first_name": randomString(),
            "last_name": randomString(),
        }    
    )

    if not resp.ok:
        print(resp.text)
    else:
        id = resp.json()['data']['id']
        student_ids[id] = []
        student_arr.append(id)
def coursesCreate():
    resp = requests.post(
        f'{courses_api}',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
        json={
            "name": randomString(),
            "description": randomString(),
            "discipline": randomString(),
            "teacher": randomString()
        }
    )

    if not resp.ok:
        print(resp.text)
    else:
        course_ids.append(resp.json()['data']['id'])
        course_ids[id] = []
        course_arr.append(id)
def studentsUpdate():
    course_id = random.choice(course_arr)
    student_id = random.choice(student_arr)
    student_ids[student_id].append(course_id)
    resp = requests.put(
        f'{students_api}/{student_id}',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
        json={
            "name": randomString(),
            "first_name": randomString(),
            "last_name": randomString(),
            "course_ids": student_ids[course_id]
        }
    )

    if not resp.ok:
        print(resp.text)
def coursesUpdate():
    course_id = random.choice(course_arr)
    student_id = random.choice(student_arr)
    course_ids[course_id].append(student_id)
    resp = requests.put(
        f'{courses_api}/{course_id}',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
        json={
            "name": randomString(),
            "description": randomString(),
            "discipline": randomString(),
            "teacher": randomString(),
            "student_ids": course_ids[course_id]
        }
    )

    if not resp.ok:
        print(resp.text)
    else:
        course_ids.append(resp.json()['data']['id'])
def studentsDelete():
    student_id = random.choice(student_arr)
    resp = requests.delete(
        f'{students_api}/{student_id}',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
    )

    if not resp.ok:
        print(resp.text)
    else:
        student_arr.remove(student_id)
        del student_ids[student_id]
def coursesDelete():
    course_id = random.choice(course_arr)
    resp = requests.delete(
        f'{courses_api}/{course_id}',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
    )

    if not resp.ok:
        print(resp.text)
    else:
        course_arr.remove(course_id)
        del course_ids[course_id]


def studentsGetAll():
    resp = requests.get(
        f'{students_api}',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
    )

    if not resp.ok:
        print(resp.text)
def coursesGetAll():
    resp = requests.get(
        f'{courses_api}',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
    )

    if not resp.ok:
        print(resp.text)
def studentsGet():
    student_id = random.choice(student_arr)
    resp = requests.get(
        f'{students_api}/{student_id}',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
    )

    if not resp.ok:
        print(resp.text)
def coursesGet():
    course_id = random.choice(course_arr)
    resp = requests.get(
        f'{courses_api}/{course_id}',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
    )

    if not resp.ok:
        print(resp.text)
def studentsInnerGet():
    student_id = random.choice(student_arr)
    resp = requests.get(
        f'{students_api}/{student_id}/courses',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
    )

    if not resp.ok:
        print(resp.text)
def coursesInnerGet():
    course_id = random.choice(course_arr)
    resp = requests.get(
        f'{courses_api}/{course_id}/students',
        headers={'x-api-key': 'T09QUyEhISBGb3VuZCBtZTopIEkgQU0gUFJPVUQgT0YgWU9V'},
    )

    if not resp.ok:
        print(resp.text)


studentQuery = [studentsGetAll, studentsGet, studentsInnerGet]
courseQuery = [coursesGetAll, coursesGet, coursesInnerGet]
studentMutation = [studentsCreate, studentsUpdate, studentsDelete]
courseMutation = [coursesCreate, coursesUpdate, coursesDelete]

def callIpStudentsInit():
    for i in range(20):
        studentsCreate()
def callIpCoursesInit():
    for i in range(20):
        coursesCreate

with ThreadPoolExecutor() as e:
    for i in range(20):
        e.submit(studentsCreate)
        e.submit(coursesCreate)

with ThreadPoolExecutor() as e:
    for i in range(10):
        e.submit(studentsUpdate)
        e.submit(coursesUpdate)

# def callIpStudentsInit():
#     for i in range(100):
#         studentsCreate()
# def callIpCoursesInit():
#     for i in range(100):
#         coursesCreate

# def callIpStudentMutate():
#     random.choice(studentMutation)()
# def callIpCourseMutate():
#     random.choice(courseMutation)()
# def callIpStudentQuery():
#     random.choice(studentQuery)()
# def callIpCourseQuery():
#     random.choice(courseQuery)()

# with ThreadPoolExecutor() as e:
#     f = e.submit(callIpStudentsInit)
#     f = e.submit(callIpCoursesInit)

# for i in range(100):
#     with ThreadPoolExecutor() as e:
#         for i in range(10):
#             e.submit(callIpStudentMutate)
#             e.submit(callIpCourseMutate)
#     with ThreadPoolExecutor() as e:
#         for i in range(100):
#             e.submit(callIpStudentQuery)
#             e.submit(callIpCourseQuery)

print(student_ids)
print(len(student_ids))
